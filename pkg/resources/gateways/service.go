/*
Copyright 2019 Banzai Cloud.

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

package gateways

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/banzaicloud/istio-operator/pkg/resources/templates"
	"github.com/banzaicloud/istio-operator/pkg/util"
)

func (r *Reconciler) service(gw string) runtime.Object {
	gwConfig := r.getGatewayConfig(gw)
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMetaWithAnnotations(gatewayName(gw), util.MergeLabels(gwConfig.ServiceLabels, labelSelector(gw)), gwConfig.ServiceAnnotations, r.Config),
		Spec: apiv1.ServiceSpec{
			Type:     r.serviceType(gw),
			Ports:    r.servicePorts(gw),
			Selector: labelSelector(gw),
		},
	}
}

func (r *Reconciler) servicePorts(gw string) []apiv1.ServicePort {
	switch gw {
	case ingress:
		ports := []apiv1.ServicePort{
			{Port: 15020, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15020), Name: "status-port", NodePort: 31460},
			{Port: 80, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(80), Name: "http2", NodePort: 31380},
			{Port: 443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(443), Name: "https", NodePort: 31390},
			{Port: 31400, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(31400), Name: "tcp", NodePort: 31400},
			{Port: 15029, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15029), Name: "https-kiali", NodePort: 31410},
			{Port: 15030, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15030), Name: "https-prom", NodePort: 31420},
			{Port: 15031, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15031), Name: "https-grafana", NodePort: 31430},
			{Port: 15032, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15032), Name: "https-tracing", NodePort: 31440},
			{Port: 15443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15443), Name: "tls", NodePort: 31450},
		}
		if util.PointerToBool(r.Config.Spec.MeshExpansion) {
			ports = append(ports, []apiv1.ServicePort{
				{Port: 15011, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15011), Name: "tcp-pilot-grpc-tls", NodePort: 31470},
				{Port: 15004, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15004), Name: "tcp-mixer-grpc-tls", NodePort: 31480},
				{Port: 8060, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(8060), Name: "tcp-citadel-grpc-tls", NodePort: 31490},
				{Port: 853, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(853), Name: "tcp-dns-tls", NodePort: 31500},
			}...)
		}
		return ports
	case egress:
		return []apiv1.ServicePort{
			{Port: 80, Name: "http2", Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(80)},
			{Port: 443, Name: "https", Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(443)},
			{Port: 15443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15443), Name: "tls"},
		}
	}
	return []apiv1.ServicePort{}
}

func (r *Reconciler) serviceType(gw string) apiv1.ServiceType {
	switch gw {
	case ingress:
		return r.Config.Spec.Gateways.IngressConfig.ServiceType
	case egress:
		return r.Config.Spec.Gateways.EgressConfig.ServiceType
	}
	return ""
}
