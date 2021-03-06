package installation

import (
	"net"
	"net/url"
	"time"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/auth/vault"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/giantswarm/api"
	"github.com/giantswarm/architect/configuration/giantswarm/desmotes"
	"github.com/giantswarm/architect/configuration/giantswarm/happa"
	"github.com/giantswarm/architect/configuration/giantswarm/passage"
	"github.com/giantswarm/architect/configuration/guest"
	"github.com/giantswarm/architect/configuration/guest/hyperkube"
	"github.com/giantswarm/architect/configuration/guest/kubectl"
	"github.com/giantswarm/architect/configuration/guest/kubernetes"
	"github.com/giantswarm/architect/configuration/monitoring"
	"github.com/giantswarm/architect/configuration/monitoring/prometheus"
	"github.com/giantswarm/architect/configuration/monitoring/testbot"
	"github.com/giantswarm/architect/configuration/provider"
	"github.com/giantswarm/architect/configuration/provider/kvm"
	"github.com/giantswarm/architect/configuration/provider/kvm/dns"
	"github.com/giantswarm/architect/configuration/provider/kvm/flannel"
	"github.com/giantswarm/architect/configuration/provider/kvm/ingress"
	"github.com/giantswarm/architect/configuration/provider/kvm/ntp"
)

var Leaseweb = configuration.Installation{
	V1: configuration.V1{
		Auth: auth.Auth{
			Vault: vault.Vault{
				Address: url.URL{
					Scheme: "https",
					Host:   "vault.g8s.fra-1.giantswarm.io:8200",
				},
				CA: vault.CA{
					TTL: 10 * 365 * 24 * time.Hour,
				},
				Certificate: vault.Certificate{
					TTL: 26 * 7 * 24 * time.Hour,
				},
				Token: vault.Token{
					TTL: 26 * 7 * 24 * time.Hour,
				},
			},
		},

		GiantSwarm: giantswarm.GiantSwarm{
			API: api.API{
				Address: url.URL{
					Scheme: "https",
					Host:   "api-g8s.giantswarm.io",
				},
			},
			Passage: passage.Passage{
				Address: url.URL{
					Scheme: "https",
					Host:   "passage-g8s.giantswarm.io",
				},
			},
			Desmotes: desmotes.Desmotes{
				Address: url.URL{
					Scheme: "https",
					Host:   "desmotes-g8s.giantswarm.io",
				},
			},
			Happa: happa.Happa{
				Address: url.URL{
					Scheme: "https",
					Host:   "happa-g8s.giantswarm.io",
				},
			},
		},

		Guest: guest.Guest{
			Hyperkube: hyperkube.Hyperkube{
				Version: hyperkube.Version,
			},
			Kubectl: kubectl.Kubectl{
				Version: kubectl.Version,
			},
			Kubernetes: kubernetes.Kubernetes{
				API: kubernetes.API{
					EndpointBase: "g8s.fra-1.giantswarm.io",
				},
			},
		},

		Monitoring: monitoring.Monitoring{
			Prometheus: prometheus.Prometheus{
				Address: url.URL{
					Scheme: "https",
					Host:   "prometheus-g8s.giantswarm.io",
				},
				RetentionPeriod: 2 * 7 * 24 * time.Hour,
			},
			Testbot: testbot.Testbot{
				Interval: 5 * time.Minute,
			},
		},

		Provider: provider.Provider{
			KVM: kvm.KVM{
				DNS: dns.DNS{
					Servers: []net.IP{
						net.ParseIP("8.8.8.8"),
						net.ParseIP("8.8.4.4"),
					},
				},
				Flannel: flannel.Flannel{
					VNIRange: flannel.Range{
						Min: 2,
						Max: 210,
					},
				},
				Ingress: ingress.Ingress{
					PortRange: ingress.PortRange{
						Min: 31000,
						Max: 31021,
					},
				},
				NTP: ntp.NTP{
					Servers: []net.IP{},
				},
			},

			Kind: provider.KVM,
		},
	},
}
