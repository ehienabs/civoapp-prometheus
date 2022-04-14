# Kubernetes Monitoring with Prometheus and Grafana

There are more than a few benefits of having insights into how our applications, and the resources that power them are performing.

We can improve our users’ experience, gain an understanding of how our services are used, help reduce mean time to resolution (MTTR) when our services run into trouble, and reduce overall downtime.

Metrics are measurements that capture data about the performance of a process over a period of time. They represent a designator or identifier where data points are updated continuously.

Our applications’ response times, error rates, uptimes, etc., are useful metrics we could use to measure the performance of our apps.  Metrics like CPU and Memory utilization, read and write operations etc., are useful for learning more about our applications’ resource usage.

## Prometheus

Prometheus is an open-source monitoring and alerting toolkit which collects and stores metrics as time series. It has a multidimensional data model which uses key/value pairs to identify data, a fast and efficient query language (PromQL), service discovery, and does not rely on distributed storage.

Using client libraries we can leverage Prometheus to monitor our services.

Client libraries allow us define internal metrics for our applications using the same programming language.

These metrics are then exposed via an HTTP endpoint which Prometheus proceeds to scrape according to a set of rules we configure. 

The Prometheus community offer several languages we can use for our app instrumentation, including Golang, Python, Java, and Ruby. There are also various third party libraries maintained by various communities including Bash, C++, Perl etc.

### Prometheus Metrics

Using Prometheus client libraries, we are able to instrument four main types of metrics in our application.

They include;

- **Counter**; When we want to know how often an event occurred we use the counter metric type.  It is a cumulative metric whose value can only be increased or reset. Examples of counter metrics include the total number of HTTP requests our service has received, the number of versions of our service etc.
- **Gauge;** When we want to take a snapshot of a metric at a point in time, we use the gauge metric type. The gauge metric type is similar to the counter in that they both take measures of the total occurrence of an event. Whereas the counter metric type measures the total occurrence of an event over a period of time, the gauge metric type takes a snapshot of the total occurrence of an event at a point in time. Gauge is used for metrics that can either go up or down. Examples of gauge metrics include current temperature, number of concurrent connections, number of online users, number of items in a queue etc.
- **Histogram;** When we want to group observations by their frequency and place them in pre-defined buckets, we use the histogram metric type.  For example, we can create buckets that specify the requests durations we want to keep track of, e.g. the Go client library uses .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10(s) as the defaults when left unspecified. So when a request is made to our service the histogram metric calculates the duration, and stores it in the appropriate bucket. That way, if we want to know how many requests took less than .005 seconds we can get that information. And so it goes for the rest of the buckets. Examples of histogram metrics include request duration, byte size, response size, etc.
- **Summary**; Summaries are similar to histograms in that they both track distributions. They are best used for measuring latencies, especially when a near accurate value is desired. Care is taken when using summaries as they can not accurately perform aggregations and tend to be expensive in terms of resources.

### Prometheus Labels

Data in Prometheus is stored as a time series. A combination of metric and labels provide an identifier for our time series data.

Labels are attributes of metrics that provide dimensionality. They are key-value pairs which enrich metrics by providing a unique identifier for each time series, enabling aggregation and filtering.

Examples of labels include; 

Instance - the name of the instance being monitored.

Handler-  the function being executed

Status_code- the returned status code. 

### Prometheus Exporters

There are services, such as those with legacy code, or third-party software for which we do not have access to their code but we would like to monitor. For these types of services, Prometheus offers exporters.

Prometheus exporters help us monitor systems we are unable to instrument metrics for. They fetch non-prometheus metrics, statistics, and other type of data, convert them into the prometheus metric format, start a server and expose these metrics at the /metric endpoint.

The Prometheus community offer and maintain a host of official exporters including;

Exporters for HTTP such as Webdriver exporter, Apache exporter, HAProxy exporter etc.

Exporters for messaging systems such as Kafka exporter, RabbitMQ exporter, Beanstalkd exporter, etc.

Exporters for Databases such as MySQL server exporter, Oracle database exporter, Redis exporter, etc.

Exporters for APIs and other monitoring systems. 

Developers are also encouraged to write their own exporters should none of the available ones meet their requirements.

## Monitoring Kubernetes resources with Prometheus

Kubernetes is a container orchestration platform. It is a portable, extensible, open-source platform for managing containerized applications. A Kubernetes cluster can manage thousands of microservices packaged and run as containers making it ideal for running services at scale.

Using Kubernetes to deploy and manage our containers we can improve our development cycle.

Prometheus integrates well with Kubernetes. Supporting service discovery, Prometheus automatically pulls metrics from newly created replicas as we scale up our services. Kubernetes and Prometheus also work with labels, which are used to select and aggregate objects for  queries, solidifying a perfect match.

### Installing Prometheus on Kubernetes

We will be using Civo managed Kubernetes service for our cluster.

Civo’s cloud-native infrastructure services are powered by Kubernetes and use the lightweight Kubernetes distribution, K3s, for superfast launch times.  

## Prerequisites

To get started, we will need the following:

- [Civo Account](https://www.civo.com/signup)
- Civo [Command Line Client](https://www.github.com/civo/cli)
- [Helm](https://helm.sh/docs/intro/install/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

After setting up the Civo command line with our API key using [the instructions in the repository](https://github.com/civo/cli#api-keys), we can create our cluster using the following command:

```bash
civo kubernetes create civo-cluster
```

and our cluster, named ‘civo-cluster’ is created 

![Annotation 2022-04-06 173836.png](Kubernetes%2060ea6/Annotation_2022-04-06_173836.png)

Kube-prometheus-stack provides a way to easily install Prometheus on Kubernetes using Helm. It is a collection of manifests that when applied to the cluster, creates an end to end monitoring stack.

The kube-prometheus-stack chart creates the following resources when applied to a cluster 

- The Prometheus Operator
- Alertmanager
- Prometheus node exporter
- Prometheus adapter for Kubernetes Metrics API
- Kube State Metrics
- Grafana

To install the Kube-Prometheus-Stack we begin by first adding the repository;

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

Then we update our repositories; 

```bash
helm repo update
```

Finally, we install our chart using this command;

```bash
helm install prom-stack prometheus-community/kube-prometheus-stack
```

Where ‘prom-stack’ is our release name.

![Annotation 2022-04-01 205916.png](Kubernetes%2060ea6/Annotation_2022-04-01_205916.png)

By running the following command

```bash
Kubectl get pods
```

We can see all the components deployed with our chart.

![Annotation 2022-04-01 210048.png](Kubernetes%2060ea6/Annotation_2022-04-01_210048.png)

Using port forwarding we can expose our Prometheus service externally with the following command;

```bash
kubectl port-forward prometheus-prom-stack-kube-prometheus-prometheus-0 9090:9090
```

Where ‘prometheus-prom-stack-kube-prometheus-prometheus-0’ is the name of our Prometheus service, and ‘9090’ is the port number.

Now we can access the user interface at [http://localhost:9090/](http://localhost:9090/graph?g0.expr=&g0.tab=1&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h) 

![Annotation 2022-04-01 211044.png](Kubernetes%2060ea6/Annotation_2022-04-01_211044.png)

Using the metrics explorer we can view the list of available metrics automatically pulled by Prometheus.

There are metrics about the API server, pods, containers, nodes, deployments, alert manager and even the Prometheus server itself.

## Grafana

Grafana is an open-source visualization and analytics software. It allows us pull metrics from a variety of sources, run queries against them, and visualize them making it easy to gain insight and make decisions about our services.

Using Grafana Dashboards, we can pull metrics from various data sources such as InfluxDB, MySQL, Datadog, etc.

Grafana integrates seamlessly with Prometheus and is deployed as one of the components when we use the Kube-Prometheus-stack to install our monitoring tools.

We can expose our Grafana service by port-forwarding using the following command;

```bash
kubectl port-forward prom-stack-grafana-6c56cdfbfb-hnlhc 3000:3000
```

Where ‘prom-stack-grafana-6c56cdfbfb-hnlhc’ is our grafana service and ‘3000’ is the port number.

Now we can access our Grafana user interface at [http://localhost:3000/](http://localhost:3000/datasources)

![Annotation 2022-04-02 014957.png](Kubernetes%2060ea6/Annotation_2022-04-02_014957.png)

To login to our service, we use ‘admin’ as the username and we can get the default password from our Grafana secret using the following command 

```bash
kubectl get secret prom-stack-grafana -o jsonpath="{.data.admin-password}" | base64
```

and we are in!

![Annotation 2022-04-02 020415.png](Kubernetes%2060ea6/Annotation_2022-04-02_020415.png)

And because we deployed our monitoring stack using Kube-Prometheus-Stack charts, our Prometheus server has already been added as a data source.

![Annotation 2022-04-02 020630.png](Kubernetes%2060ea6/Annotation_2022-04-02_020630.png)
