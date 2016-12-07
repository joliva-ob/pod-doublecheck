package kubernetes


import (

  "flag"

  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/pkg/api/v1"
  "k8s.io/client-go/tools/clientcmd"
  "github.com/joliva-ob/pod-doublecheck/config"

)


var (
  kubeconfig = flag.String("kubeconfig", "/Users/joan/.kube/config", "absolute path to the kubeconfig file")
)


func GetPodsMap() map[string]bool {

    podsMap := make(map[string]bool)  // k: pod name v: found status, start from true
    flag.Parse()

    // uses the current context in kubeconfig
    kconfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
      panic(err.Error())
    }

    // creates the clientset
    clientset, err := kubernetes.NewForConfig(kconfig)
    if err != nil {
      panic(err.Error())
    }

    pods, err := clientset.Core().Pods("pro").List(v1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }
    for _, p := range pods.Items {
        podsMap[p.Status.ContainerStatuses[0].Name] = true
        config.Log.Debugf("pod name: %v", p.Status.ContainerStatuses[0].Name)
    }

    return podsMap
}