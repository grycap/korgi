/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	esupvgrycapv1 "es.upv.grycap/korgi/api/v1"
)

const (
	korgiJobSchedulerField = ".spec.KorgiJobScheduler"
)

var korgiJobSchedulerHasChanged bool = false
var korgiJobSchedulerVersion string = ""
var korgiJobSchedulerGPU string = ""

// KorgiJobReconciler reconciles a KorgiJob object
type KorgiJobReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	PrevState map[string]esupvgrycapv1.KorgiJob
}

/* NewKorgiJobReconciler creates a new KorgiJobReconciler
func NewKorgiJobReconciler(client client.Client, scheme *runtime.Scheme) *KorgiJobReconciler {
	return &KorgiJobReconciler{
		Client:    client,
		Scheme:    scheme,
		PrevState: make(map[string]esupvgrycapv1.KorgiJob),
	}
}
*/

//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=jobs/status,verbs=get
//+kubebuilder:rbac:groups="",resources=korgijobschedulers,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KorgiJob object against the actual cluster state, and then perform
// operations to make the cluster state reflect the state specified by the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *KorgiJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	log.Info("Reconciling KorgiJob: ", "namespace", req.Namespace, "name ", req.Name)

	// Get the current state of the object from the API server
	var korgiJob = &esupvgrycapv1.KorgiJob{}
	if err := r.Client.Get(ctx, req.NamespacedName, korgiJob); err != nil {
		log.Error(err, "Unable to fetch KorgiJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("KorgiJob: " + korgiJob.GetName())

	// Check if we have a previous state for this object
	if _, ok := r.PrevState[req.NamespacedName.String()]; !ok {
		// If not, create a new entry in the map
		r.PrevState[req.NamespacedName.String()] = *korgiJob.DeepCopy()
	}
	// If ok, we have a previous state, so we can compare it with the current one
	prevState := r.PrevState[req.NamespacedName.String()]
	// Compare the previous and current states to determine what changed using reflect.DeepEqual
	if !reflect.DeepEqual(prevState, *korgiJob) {
		// print the diff
		log.Info("KorgiJob changed")
		log.Info("Diff: ", "diff", cmp.Diff(prevState, *korgiJob))
	}
	//Update previous state
	r.PrevState[req.NamespacedName.String()] = *korgiJob.DeepCopy()

	start := time.Now()

	korgiJobSchedulerName := korgiJob.Spec.KorgiJobScheduler
	foundKorgiJobScheduler := &esupvgrycapv1.KorgiJobScheduler{}
	err := r.Get(ctx, types.NamespacedName{Name: korgiJobSchedulerName, Namespace: korgiJob.Namespace}, foundKorgiJobScheduler)
	for err != nil {
		if time.Since(start) >= 2*time.Minute { //TO DO: make it configurable (hardcoded TimeOut duration)
			log.Info("KorgiJob TIMEOUT")
			// Update KorgiJob status
			korgiJob.Status.Status = esupvgrycapv1.KorgiJobFailed
			if err2 := r.Client.Status().Update(ctx, korgiJob); err2 != nil {
				log.Error(err, "KorgiJobScheduler couldn't be retrieved")
				log.Error(err2, "Status update failed (to failed due to timeout )")
				return ctrl.Result{}, err2
			}
			break
		}
		err = r.Get(ctx, types.NamespacedName{Name: korgiJobSchedulerName, Namespace: korgiJob.Namespace}, foundKorgiJobScheduler)
	}
	// Hash the data in some way, or just use the version of the Object
	log.Info(korgiJobSchedulerVersion + " " + foundKorgiJobScheduler.ResourceVersion)
	log.Info("korgiJob name for KJS " + foundKorgiJobScheduler.Spec.KorgiJobName)
	log.Info(korgiJob.GetName())
	if korgiJobSchedulerVersion != foundKorgiJobScheduler.ResourceVersion && foundKorgiJobScheduler.Spec.KorgiJobName == korgiJob.GetName() {
		log.Info("New KorgiJobScheduler!!")
		korgiJobSchedulerVersion = foundKorgiJobScheduler.ResourceVersion
		korgiJobSchedulerGPU = foundKorgiJobScheduler.Spec.GPUResources
		korgiJobSchedulerHasChanged = true
	}

	switch korgiJob.GetStatus() {
	case "":

		if korgiJobSchedulerHasChanged {
			log.Info("KorgiJob.Status = EMPTY")
			korgiJobSchedulerHasChanged = false
			// Update KorgiJob status
			korgiJob.Status.Status = esupvgrycapv1.KorgiJobPending
			if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
				log.Error(err, "Status update failed (to pending ) 010")
				return ctrl.Result{}, err
			}
		}
		// Check creation time
		// if > x hours	-> change status to Failed
		// else 		-> wait for a KorgiJobScheduler
	case esupvgrycapv1.KorgiJobPending:
		// Create associated job
		log.Info("KorgiJob.Status = PENDING")
		log.Info("Creating job from KorgiJob")
		jobName := fmt.Sprintf("%s-%d", korgiJob.Name+"-subjob", time.Now().Unix())
		annotations := make(map[string]string)
		annotations["KorgiJobSchedulerVersion"] = korgiJobSchedulerVersion
		ttlSecondsAfterFinished := int32(3600) //TO DO: make it configurable
		job := batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:        jobName,
				Namespace:   korgiJob.Namespace,
				Labels:      korgiJob.Labels,
				Annotations: annotations,
			},

			Spec: batchv1.JobSpec{
				TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							corev1.Container{
								Name:    korgiJob.Name + "-container",
								Image:   korgiJob.Spec.Image,
								Command: korgiJob.Spec.Command,
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceName(korgiJobSchedulerGPU): resource.MustParse("1"),
									},
								},
							},
						},
						RestartPolicy: "Never",
					},
				},
			},
		}
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, &job, func() error {
			log.Info("Job created", "job", job.Name)
			ctrl.SetControllerReference(korgiJob, &job, r.Scheme)
			// Update KorgiJob status
			korgiJob.Status.Status = esupvgrycapv1.KorgiJobRunning
			if err2 := r.Client.Status().Update(ctx, korgiJob); err2 != nil {
				log.Error(err2, "Status update failed (from pending to running) 011")
				return err2
			}
			return nil
		}); err != nil {
			return ctrl.Result{}, err
		}

	case esupvgrycapv1.KorgiJobRunning:
		// Check Job Status
		// if Active 	-> wait
		// if Succeded 	-> change status to Completed
		// if Failed  	-> change status to Failed
		log.Info("KorgiJob.Status = RUNNING")
		if korgiJobSchedulerHasChanged {
			korgiJobSchedulerHasChanged = false
			korgiJob.Status.Status = esupvgrycapv1.KorgiJobRescheduling
			if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
				log.Error(err, "Status update failed 012b")
				return ctrl.Result{}, err
			}
		} else {
			var childJobs batchv1.JobList
			if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace)); err != nil {
				log.Error(err, "Unable to list child Jobs")
				return ctrl.Result{}, err
			}
			if childJobs.Size() > 0 {
				active := childJobs.Items[0].Status.Active
				succeeded := childJobs.Items[0].Status.Succeeded
				failed := childJobs.Items[0].Status.Failed
				log.Info("", "Active: ", active, "Succeeded: ", succeeded, "Failed: ", failed)
				if !(active == 0 && failed == 0 && succeeded == 0) {
					if failed > 0 {
						korgiJob.Status.Status = esupvgrycapv1.KorgiJobFailed
					}
					if succeeded > 0 {
						korgiJob.Status.Status = esupvgrycapv1.KorgiJobCompleted
					}
					if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
						log.Error(err, "Status update failed 012")
						return ctrl.Result{}, err
					}
				}
			}
		}
	case esupvgrycapv1.KorgiJobRescheduling:
		log.Info("KorgiJob.Status = RESCHEDULING")
		// Delete any child jobs
		var childJobs batchv1.JobList
		if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace)); err != nil {
			log.Error(err, "Unable to list child Jobs")
			return ctrl.Result{}, err
		}
		for _, job := range childJobs.Items {
			if err := r.Delete(ctx, &job); err != nil {
				log.Error(err, "Unable to delete child Job")
				return ctrl.Result{}, err
			}
			log.Info("ChildJob deleted", "job", job.GetName())
		}
		// Change status to Pending
		korgiJob.Status.Status = esupvgrycapv1.KorgiJobPending
		if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
			log.Error(err, "Status update failed (from resch. to pending) 013")
			return ctrl.Result{}, err
		}
	case esupvgrycapv1.KorgiJobCompleted:
		// --
		log.Info("KorgiJob.Status = COMPLETED")

	case esupvgrycapv1.KorgiJobFailed:
		log.Info("KorgiJob.Status = FAILED")
		if korgiJobSchedulerHasChanged {
			// Update KorgiJob status
			korgiJob.Status.Status = esupvgrycapv1.KorgiJobPending
			if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
				log.Error(err, "Status update failed (to pending ) 010")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *KorgiJobReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &esupvgrycapv1.KorgiJob{}, korgiJobSchedulerField, func(rawObj client.Object) []string {
		// Extract the KorgiJobScheduler name from the KorgiJob Spec, if one is provided
		korgiJob := rawObj.(*esupvgrycapv1.KorgiJob)
		if korgiJob.Spec.KorgiJobScheduler == "" {
			return nil
		}
		return []string{korgiJob.Spec.KorgiJobScheduler}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&esupvgrycapv1.KorgiJob{}).
		Owns(&batchv1.Job{}).
		Watches(
			&source.Kind{Type: &esupvgrycapv1.KorgiJobScheduler{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForKorgiJobScheduler),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *KorgiJobReconciler) findObjectsForKorgiJobScheduler(korgiJobScheduler client.Object) []reconcile.Request {
	attachedKorgiJobs := &esupvgrycapv1.KorgiJobList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(korgiJobSchedulerField, korgiJobScheduler.GetName()),
		Namespace:     korgiJobScheduler.GetNamespace(),
	}
	err := r.List(context.TODO(), attachedKorgiJobs, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(attachedKorgiJobs.Items))
	for i, item := range attachedKorgiJobs.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}
	return requests
}
