/*
Copyright 2022.

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

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	korgiv1 "korgi.grycap.upv.es/korgi-operator/api/v1"
)

// KorgiJobReconciler reconciles a KorgiJob object
type KorgiJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=korgi.korgi.grycap.upv.es,resources=korgijobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=korgi.korgi.grycap.upv.es,resources=korgijobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=korgi.korgi.grycap.upv.es,resources=korgijobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KorgiJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *KorgiJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	korgiJob := &korgiv1.KorgiJob{}
	if err := r.Client.Get(ctx, req.NamespacedName, korgiJob); err != nil {
		log.Error(err, "Unable to fetch KorgiJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	switch korgiJob.GetStatus() {
	case "":
		log.Info(("Creating job from KorgiJob"))
		job := batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      korgiJob.Name + "-subjob",
				Namespace: korgiJob.Namespace,
				Labels:    korgiJob.Labels,
			},

			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							corev1.Container{
								Name:    korgiJob.Name + "-container",
								Image:   korgiJob.Spec.Image,
								Command: korgiJob.Spec.Command,
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										"nvidia.com/gpu": resource.MustParse("1"),
									},
								},
							},
						},
						RestartPolicy: "Never",
					},
				},
			},
		}
		log.Info("Job created", "job", job.Name)
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, &job, func() error {
			if err := ctrl.SetControllerReference(korgiJob, &job, r.Scheme); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return ctrl.Result{}, err
		}
		korgiJob.Status.Status = korgiv1.KorgiJobRunning
		if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
			log.Error(err, "Status update failed")
			return ctrl.Result{}, err
		}
		log.Info("Changing status to pending")
		korgiJob.Status.Status = korgiv1.KorgiJobPending
		if err := r.Client.Status().Update(ctx, korgiJob); err != nil {
			log.Error(err, "Status update failed")
			return ctrl.Result{}, err
		}
	case korgiv1.KorgiJobPending:
		var childJobs batchv1.JobList
		if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace)); err != nil {
			log.Error(err, "Unable to list child Jobs")
			return ctrl.Result{}, err
		}
		log.Info(childJobs.Items[0].String())
	case korgiv1.KorgiJobRunning:
		log.Info("Checking subjob status")
	case korgiv1.KorgiJobCompleted:
		log.Info("KorgiJob completed")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KorgiJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&korgiv1.KorgiJob{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
