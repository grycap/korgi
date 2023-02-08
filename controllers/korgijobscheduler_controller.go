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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	// TODO(user): your logic here
	switch korgiJobScheduler.GetStatus() {
	case "":
		log.Info("..... KorgiJobScheduler Controller Starts .....")
		log.Info("KorgiJob Name = " + korgiJobScheduler.Spec.KorgiJobName)
		log.Info("KorgiJob GPU Resources = " + korgiJobScheduler.Spec.GPURecources)
		// Check if the KorgiJob already exists
		found := &esupvgrycapv1.KorgiJob{}
		err := r.Get(context.TODO(), types.NamespacedName{Name: korgiJobScheduler.Spec.KorgiJobName, Namespace: req.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("KorgiJob " + korgiJobScheduler.Spec.KorgiJobName + " not found.")
			// if creation time > X hours KJS.Status = FAILED
			log.Info("korgiJobScheduler Status: " + korgiJobScheduler.GetStatus())
			log.Info("korgiJobScheduler Creation Time: " + korgiJobScheduler.ObjectMeta.CreationTimestamp.String())
			log.Info("Time now: " + time.Now().String())
			//hour := 3600000
			var hour int64 = 5000
			lifetime := time.Now().UnixMilli() - korgiJobScheduler.ObjectMeta.CreationTimestamp.Time.UnixMilli()
			if lifetime > hour {
				// if creation time > X hours KJS.Status = FAILED
				log.Info("HAN PASADO 5 SEGUNDOS!")
				// TODO never entering here
			} else {
				// This is an ugly patch, problem, it generates a lot of not wanted INFO (loop until time is reached)
				// r.Reconcile(ctx, req)
			}
		} else {
			log.Info("KorgiJob " + korgiJobScheduler.Spec.KorgiJobName + " has been found.")
			// Pass GPU information to korgijob
			// Change KJS.Status to Running
		}

	case esupvgrycapv1.KorgiJobSchedulerRunning:
		// Check KJ Status
		switch "" /* KJ.Status */ {
		case "":
			// Change korgijob status to pending
		case esupvgrycapv1.KorgiJobRunning:
			// Change KJ.Status to Rescheduling
			// Change KJS.Status to Completed
		case esupvgrycapv1.KorgiJobCompleted:
			// Change KJS.Status to Completed
		case esupvgrycapv1.KorgiJobFailed:
			// Change KJ.Status to Rescheduling
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KorgiJobSchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&esupvgrycapv1.KorgiJobScheduler{}).
		Complete(r)
}
