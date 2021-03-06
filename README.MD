# OpenStack Meetup Barcelona - July 2018 #

## Operators Demo ##

The environment for these demos is an OpenShift Origin Cluster v3.10 deployed using the following command (remember to have a 'origin-clients' versioned 3.10 or higher for this to work):

~~~sh
$ oc cluster up --enable=service-catalog,router,registry,web-console,persistent-volumes,rhel-imagestreams,automation-service-broker
~~~

At this point you'll need to enable `admin` to have cluster-admin role, hence, execute:

~~~sh
oc login -u system:admin
oc adm policy add-cluster-role-to-user cluster-admin admin
~~~

For the second demo golang and dep must be installed on your computer.

### Demo 1 - ETCD Operator  ###

In this first demo, the etcd operator will be deployed on our OpenShift Cluster. Once the operator is deployed we can create etcd clusters, definitions, and the operator will take care of those and make the required actions to ensure the cluster state is what we have asked for.

1. Clone the etcd operator repository
    ~~~sh
    $ git clone https://github.com/coreos/etcd-operator
    ~~~
2. Create a new project where the operator will be deployed
    ~~~sh
    $ oc new-project etcdoperator
    ~~~
3. The operator needs some permissions, the create_role.sh script will take care of that
    ~~~sh
    $ etcd-operator/example/rbac/create_role.sh --role-name=etcd-operator --namespace=etcdoperator --role-binding-name=etcd-operator
    ~~~
4. Create the operator's deployment in order to deploy the etcd operator
    ~~~sh
    $ oc create -f etcd-operator/example/deployment.yaml
    ~~~
5. The operator has been programmed to create the CRD by itself
    ~~~sh
    $ oc get crd
    ~~~
6. Create an etcd cluster definition
    ~~~sh
    $ oc create -f etcd-operator/example/example-etcd-cluster.yaml
    ~~~
7. Check the etcd cluster creation
    ~~~sh
    $ oc get pods -w
    ~~~
8. Try to connect to the cluster and insert a key in etcd
    ~~~sh
    $ docker run -ti --rm -e ETCDCTL_API=3 -e ETCDCTL_ENDPOINTS=http://$(oc get svc example-etcd-cluster-client -o jsonpath="{.spec.clusterIP}" -n etcdoperator):$(oc get svc example-etcd-cluster-client -o jsonpath="{.spec.ports[0].port}") centos:7 /bin/bash
    $ yum install -y etcd
    $ etcdctl put ostack meetup
    ~~~
9. Edit the cluster definition to scale it up
    ~~~sh
    $ oc patch etcdclusters example-etcd-cluster -p '{"spec":{"size":5}}' --type='merge'
    $ oc get pods -w
    ~~~
10. That's it!

### Demo 2 - Build your own operator  ###

In this demo an operator will be created using the Operator Framework SDK. We will be deploying a Python application (<https://quay.io/mavazque/hello-api>, the application itself is in files folder) that exposes an API that returns "Hello World".

1. Installing the Operator Framework SDK
    ~~~sh
    $ go get github.com/operator-framework/operator-sdk
    $ cd $GOPATH/src/github.com/operator-framework/operator-sdk
    $ git checkout master
    $ make dep
    $ make install
    ~~~
2. Initialize your operator's project
    ~~~sh
    $ mkdir -p $GOPATH/src/github.com/<user>/
    $ cd $_
    $ $GOPATH/bin/operator-sdk new <operator-name> --api-version=<your-crd-api-group>/v1alpha1 --kind=<your-crd-object-kind>
    ~~~
3. Modify your Operator types (example in files/types.go)
    ~~~sh
    $ vim $GOPATH/src/github.com/<operator-name>/pkg/apis/<api-group>/v1alpha1/types.go
    ~~~
4. Regenerate some code after modifying types.go
    ~~~sh
    $ $GOPATH/bin/operator-sdk generate k8s
    ~~~
5. Create your operator's business logic on handler.go (example in files/handler.go)
    ~~~sh
    $ vim $GOPATH/src/github.com/<operator-name>/pkg/stub/handler.go
    ~~~
6. Build your operator
    ~~~sh
    $GOPATH/bin/operator-sdk build quay.io/<user>/<operator-image-name>:<operator-image-tag>
    ~~~
7. Push your operator to the registry
    ~~~sh
    docker push quay.io/<user>/<operator-image-name>:<operator-image-tag>
    ~~~
8. Deploy your operator
    ~~~sh
    $ oc new-project helloworld-operator
    $ oc create -f /project/path/deploy/rbac.yaml
    $ oc create -f /project/path/deploy/crd.yaml
    $ oc create -f /project/path/deploy/operator.yaml
    ~~~
9. Create a Python API HelloWorld definition
    ~~~sh
    $ oc create -f /project/path/deploy/cr.yaml
    ~~~
10. Verify the deployment
    ~~~sh
    $ oc get pods
    $ oc get svc
    $ curl <svc-ip>:<svc-port>
    ~~~
10. That's it!
