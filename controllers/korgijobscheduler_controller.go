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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	esupvgrycapv1 "es.upv.grycap/korgi/api/v1"
)

// KorgiJobSchedulerReconciler reconciles a KorgiJobScheduler object
type KorgiJobSchedulerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobschedulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobschedulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.upv.grycap,resources=korgijobschedulers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KorgiJobScheduler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *KorgiJobSchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	korgiJobScheduler := &esupvgrycapv1.KorgiJobScheduler{}
	if err := r.Client.Get(ctx, req.NamespacedName, korgiJobScheduler); err != nil {
		log.Error(err, "Unable to fetch KorgiJobScheduler")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	/*
		foundKJ := &esupvgrycapv1.KorgiJob{}
		// Check if the KorgiJob already exists
		errFoundKJ := r.Get(ctx, types.NamespacedName{Name: korgiJobScheduler.Spec.KorgiJobName, Namespace: req.Namespace}, foundKJ)
		// TODO(user): your logic here
		switch korgiJobScheduler.GetStatus() {
		case "":
			log.Info("..... KorgiJobScheduler Controller Starts .....")
			log.Info("KorgiJob Name = " + korgiJobScheduler.Spec.KorgiJobName)
			log.Info("KorgiJob GPU Resources = " + korgiJobScheduler.Spec.GPUResources)

			if errFoundKJ != nil && errors.IsNotFound(errFoundKJ) {
				log.Info("KorgiJob " + korgiJobScheduler.Spec.KorgiJobName + " not found.")
			} else {
				log.Info("KorgiJob " + korgiJobScheduler.Spec.KorgiJobName + " has been found.")
				// Change KJS.Status to Running
				korgiJobScheduler.Status.Status = esupvgrycapv1.KorgiJobSchedulerRunning
				if err := r.Client.Status().Update(ctx, korgiJobScheduler); err != nil {
					log.Error(err, "Status update failed 001")
					return ctrl.Result{}, err
				}
			}

		case esupvgrycapv1.KorgiJobSchedulerRunning:
			// Check KJ Status
			log.Info("..... KorgiJobScheduler Controller -> KJS Running .....")
			switch foundKJ.Status.Status {
			case "":
				// Pass GPU information to korgijob
				foundKJ.Spec.GPUInfo = korgiJobScheduler.Spec.GPUResources
				/*
					// not needed anymore, status didn't change
					if err := r.Client.Status().Update(ctx, foundKJ); err != nil {
						log.Error(err, "Status update failed 002")
						return ctrl.Result{}, err
					}
				*
				log.Info("KorgiJob named (" + foundKJ.Name + ") GPU INFO updated: " + foundKJ.Spec.GPUInfo)
				// Change korgijob status to pending
				/*
					log.Info("KJ.Status = '', change korgijob status to pending...")
					foundKJ.Status.Status = esupvgrycapv1.KorgiJobPending
					if err := r.Client.Status().Update(ctx, foundKJ); err != nil {
						log.Error(err, "Status update failed")
						return ctrl.Result{}, err
					}
				*
				// Change KJS.Status to Completed
				korgiJobScheduler.Status.Status = esupvgrycapv1.KorgiJobSchedulerCompleted
				if err := r.Client.Status().Update(ctx, korgiJobScheduler); err != nil {
					log.Error(err, "Status update failed 003")
					return ctrl.Result{}, err
				}
			case esupvgrycapv1.KorgiJobRunning:
				// Pass GPU information to korgijob
				foundKJ.Spec.GPUInfo = korgiJobScheduler.Spec.GPUResources
				/*
					// not needed anymore, status didn't change
					if err := r.Client.Status().Update(ctx, foundKJ); err != nil {
						log.Error(err, "Status update failed 004")
						return ctrl.Result{}, err
					}
				*
				log.Info("KorgiJob named (" + foundKJ.Name + ") GPU INFO updated: " + foundKJ.Spec.GPUInfo)
				// Change KJ.Status to Rescheduling
				log.Info("KJ.Status = RUNNING, change korgijob status to rescheduling...")
				foundKJ.Status.Status = esupvgrycapv1.KorgiJobRescheduling
				if err := r.Client.Status().Update(ctx, foundKJ); err != nil {
					log.Error(err, "Status update failed 005")
					return ctrl.Result{}, err
				}
				// Change KJS.Status to Completed
				korgiJobScheduler.Status.Status = esupvgrycapv1.KorgiJobSchedulerCompleted
				if err := r.Client.Status().Update(ctx, korgiJobScheduler); err != nil {
					log.Error(err, "Status update failed 006")
					return ctrl.Result{}, err
				}
			case esupvgrycapv1.KorgiJobCompleted:
				log.Info("KJ.Status = Completed")
				// Change KJS.Status to Completed
				log.Info("KJ is in COMPLETED status, KJS would not be applied")
				// Change KJS.Status to Completed
				korgiJobScheduler.Status.Status = esupvgrycapv1.KorgiJobSchedulerCompleted
				if err := r.Client.Status().Update(ctx, korgiJobScheduler); err != nil {
					log.Error(err, "Status update failed 007")
					return ctrl.Result{}, err
				}
			case esupvgrycapv1.KorgiJobFailed:
				log.Info("KJ.Status = Failed")
				log.Info("KJ is in FAILED status, KJS cannot be applied")
				// Change KJS.Status to Completed
				korgiJobScheduler.Status.Status = esupvgrycapv1.KorgiJobSchedulerCompleted
				if err := r.Client.Status().Update(ctx, korgiJobScheduler); err != nil {
					log.Error(err, "Status update failed 008")
					return ctrl.Result{}, err
				}
			}
		}
	*/
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KorgiJobSchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&esupvgrycapv1.KorgiJobScheduler{}).
		Complete(r)
}
