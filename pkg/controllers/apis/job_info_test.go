/*
Copyright 2019 The Volcano Authors.

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

package apis

import (
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	vkbatchv1 "volcano.sh/volcano/pkg/apis/batch/v1alpha1"
)

func TestAddPod(t *testing.T) {
	namespace := "test"
	name := "pod1"

	testCases := []struct {
		Name        string
		jobinfo     JobInfo
		pod         *v1.Pod
		ExpectValue bool
		ExpectErr   string
	}{
		{
			Name: "AddPod",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID(fmt.Sprintf("%v-%v", namespace, name)),
					Name:      name,
					Namespace: namespace,
					Labels:    nil,
					Annotations: map[string]string{vkbatchv1.JobNameKey: "job1",
						vkbatchv1.JobVersion:  "0",
						vkbatchv1.TaskSpecKey: "task1"},
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
			jobinfo: JobInfo{
				Pods: make(map[string]map[string]*v1.Pod),
			},
			ExpectValue: true,
			ExpectErr:   "duplicated pod",
		},
	}

	for i, testcase := range testCases {
		err := testcase.jobinfo.AddPod(testcase.pod)
		if err != nil {
			t.Fatalf("AddPod() error: %v", err)
		}

		if _, ok := testcase.jobinfo.Pods["task1"][testcase.pod.Name]; ok != testcase.ExpectValue {
			t.Errorf("case %d (%s): expected: %v, got %v ", i, testcase.Name, testcase.ExpectValue, ok)
		}

		err = testcase.jobinfo.AddPod(testcase.pod)

		if err == nil {
			t.Errorf("case %d (%s): expected: %v, got %v ", i, testcase.Name, testcase.ExpectErr, nil)
		}
	}

}

func TestDeletePod(t *testing.T) {
	namespace := "test"
	name := "pod1"

	testCases := []struct {
		Name        string
		jobinfo     JobInfo
		pod         *v1.Pod
		ExpectValue bool
	}{
		{
			Name: "DeletePod",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID(fmt.Sprintf("%v-%v", namespace, name)),
					Name:      name,
					Namespace: namespace,
					Labels:    nil,
					Annotations: map[string]string{vkbatchv1.JobNameKey: "job1",
						vkbatchv1.JobVersion:  "0",
						vkbatchv1.TaskSpecKey: "task1"},
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
			jobinfo: JobInfo{
				Pods: make(map[string]map[string]*v1.Pod),
			},
			ExpectValue: false,
		},
	}

	for i, testcase := range testCases {

		testcase.jobinfo.Pods["task1"] = make(map[string]*v1.Pod)
		testcase.jobinfo.Pods["task1"][testcase.pod.Name] = testcase.pod

		err := testcase.jobinfo.DeletePod(testcase.pod)
		if err != nil {
			t.Fatalf("DeletePod() error: %v", err)
		}
		if _, ok := testcase.jobinfo.Pods["task1"][testcase.pod.Name]; ok != testcase.ExpectValue {
			t.Errorf("case %d (%s): expected: %v, got %v ", i, testcase.Name, testcase.ExpectValue, ok)
		}
	}
}

func TestUpdatePod(t *testing.T) {
	namespace := "test"
	name := "pod1"

	testCases := []struct {
		Name        string
		jobinfo     JobInfo
		oldpod      *v1.Pod
		newpod      *v1.Pod
		ExpectValue v1.PodPhase
	}{
		{
			Name: "UpdatePod",
			oldpod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID(fmt.Sprintf("%v-%v", namespace, name)),
					Name:      name,
					Namespace: namespace,
					Labels:    nil,
					Annotations: map[string]string{vkbatchv1.JobNameKey: "job1",
						vkbatchv1.JobVersion:  "0",
						vkbatchv1.TaskSpecKey: "task1"},
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
			newpod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID(fmt.Sprintf("%v-%v", namespace, name)),
					Name:      name,
					Namespace: namespace,
					Labels:    nil,
					Annotations: map[string]string{vkbatchv1.JobNameKey: "job1",
						vkbatchv1.JobVersion:  "0",
						vkbatchv1.TaskSpecKey: "task1"},
				},
				Status: v1.PodStatus{
					Phase: v1.PodSucceeded,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
			jobinfo: JobInfo{
				Pods: make(map[string]map[string]*v1.Pod),
			},
			ExpectValue: v1.PodSucceeded,
		},
	}

	for i, testcase := range testCases {

		testcase.jobinfo.Pods["task1"] = make(map[string]*v1.Pod)
		testcase.jobinfo.Pods["task1"][testcase.oldpod.Name] = testcase.oldpod

		err := testcase.jobinfo.UpdatePod(testcase.newpod)
		if err != nil {
			t.Fatalf("UpdatePod() error: %v", err)
		}
		if val, ok := testcase.jobinfo.Pods["task1"][testcase.newpod.Name]; ok != true {
			t.Errorf("case %d (%s): expected: %v, got %v ", i, testcase.Name, true, ok)
		} else if val.Status.Phase != v1.PodSucceeded {
			t.Errorf("case %d (%s): expected: %v, got %v ", i, testcase.Name, testcase.ExpectValue, val.Status.Phase)
		}
	}
}
