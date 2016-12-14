package kubernetes


import (

  "flag"
  "strconv"

  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/pkg/api/v1"
  "k8s.io/client-go/tools/clientcmd"

  "github.com/joliva-ob/pod-doublecheck/config"
  "github.com/joliva-ob/pod-doublecheck/handler"

)


var (
  kubeconfig = flag.String("kubeconfig", "kube/config", "absolute path to the kubeconfig file")
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

    pods, err := clientset.Core().Pods(config.Configuration["ENV"].(string)).List(v1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }
    for _, p := range pods.Items {
        podsMap[p.Status.ContainerStatuses[0].Name] = true
        //config.Log.Debugf("pod name: %v", p.Status.ContainerStatuses[0].Name)
    }
    handler.AddMetric("Kubernetes pods", int64(len(pods.Items)), 300) // 300 Max number of pods allowed
    config.Log.Infof(strconv.Itoa(len(pods.Items)) + " kubernetes pods found.")

    return podsMap
}