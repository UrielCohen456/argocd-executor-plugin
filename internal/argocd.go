package argocd



// output, err := exec.Command("argocd").Output()
// if err!=nil {
// 		fmt.Println(err.Error())
// }
// fmt.Println(string(output))

// pods, err := clientSet.CoreV1().Pods(ctx.Value("namespace").(string)).List(ctx, metav1.ListOptions{})
// if err != nil {
// 	panic(err.Error())
// }
// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

// _, err = clientSet.CoreV1().Pods("default").Get(ctx, "example-xxxxx", metav1.GetOptions{})
// if errors.IsNotFound(err) {
// 	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
// } else if err != nil {
// 	panic(err.Error())
// } else {
// 	fmt.Printf("Found example-xxxxx pod in default namespace\n")
// }