---
title: "Building an Overkill Home Server"
date: 2023-01-07T07:07:21+01:00
showDate: true
mathjax: true
justify: true
tags: ["devops", "kubernetes", "docker", "server"]
---

![](/content/posts/build_an_overkill_server/server-setup.png)

For future side projects and to learn Kubernetes more interactively I decided
to build a server at home. I've been thinking about it a lot
and then last summer I finally scrambled together all the components to build it.

Because I don't only want the server to be able to host websites
but also GameServers I opted for a more powerful build.

# Kubernetes, what and why?

My whole plan for this server was to use it for hosting different applications
and websites. This includes my Discord bots
([Lecturfier](https://github.com/markbeep/lecturfier), [Cup](https://github.com/markbeep/cup),
etc.) but also websites I like to hack at on the side. For example, my still WIP
[Wenjim](https://wenjim.markc.su/) website to check when the best time to go to ASVZ
is. In the future, I was also planning to upgrade this blog to a dynamic website.

For this setup I gave myself three options for how I could set up the whole server:

1. Install everything as normal and get into a package and dependency hell.
2. Use Docker for every application and run it all with Docker Compose. This would
   be simple but require me to Dockerize a few applications that aren't yet.
3. Do 2. and set up a single node Kubernetes cluster.

You can already tell, of course, I went with the Kubernetes cluster. To explain
Kubernetes very briefly. It is a server tool to orchestrate applications across
multiple nodes (servers). So instead of having to install an application on multiple
machines and handling the load balancer somewhere manually, you can simply set up
Kubernetes on all your nodes and then the distribution will be handled automatically
with an easy way to add a load balancer. One of your nodes goes down or is overused?
No worries, Kubernetes will move your applications to a different working node.

So all this talk about multiple nodes, but I'm running it on a single node. Why?
Another upside of Kubernetes is that you simply feed it a Dockerfile of your
app and a YAML file describing what it should do and it handles the rest. If you know
docker-compose, it's basically that but on steroids. This allows me to have a very
organized way to handle all my different apps and Kubernetes always makes sure they
stay running if they ever crash or I restart my server.

Another reason for choosing Kubernetes is that the
[VIS](https://vis.ethz.ch/en/) backend is also a big Kubernetes cluster.
So to learn how it all works I thought it would be best to also use it on my server.

Instead of installing the full Kubernetes (K8s) I opted for a more lightweight
install called [K3s](https://k3s.io/). I picked it because it was easier to
set up and more suited for single-node clusters. In the end, it works
the same as a normal Kubernetes setup.

After installing Kubernetes my task for all my applications was now to Dockerize
them and make them able to be run on Kubernetes. Luckily it is quite easy to
dockerize apps that simply execute python or node code.

# Building Docker Images Automatically

Now Kubernetes is a nice tool to orchestrate it all, but how do I add my
nice Docker images to it all? One way would be to do some hacky local shenanigans
to upload images to my server and have Kubernetes use them. This is a pretty annoying
workflow though. I'd have to constantly connect to my Kubernetes cluster to upload
the images and run them there. This is also not scalable to multiple nodes because
I'd have to make sure every node has access to the local Docker image.

Because of this hassle, I decided to simply host all my images on
[Docker Hub](https://hub.docker.com/). A free way to store all your Docker images
and then easily download them on any node of your cluster. Now all my Github
repositories I want to be built to Docker Hub I simply have a Github action
([example](https://github.com/markbeep/Wenjim/blob/master/.github/workflows/docker-prod.yaml))
that then automatically builds my Dockerfiles into images and uploads them to
Docker Hub. So after every push, I'll now have the most up-to-date image available
on any node (and you do as well). [This](https://github.com/markbeep/Wenjim/blob/master/.github/workflows/docker-staging.yaml)
is for example the Github action for one of my sites with a backend and frontend that
I usually just copy-paste and modify slightly for new projects.

So that's cool, but don't I still have to manually access my server to create the
Kubernetes config and restart the applications? No, that's where ArgoCD comes
into play.
https://github.com/markbeep/Wenjim/blob/master/.github/workflows/docker-staging.yaml

# ArgoCD

ArgoCD is a way to store all your deployment configs inside a Git repository
and it makes sure that the Kubernetes setup matches exactly what is on Git.
This is perfect. It allows me to store how applications should be started and
what peristent volumes and ingresses should exist without even having to access the
server. I open my Kubernetes repo, throw in another file for my new deployment,
push and then the deployment automatically starts up.

Another bonus of storing deployments in this way is that in case my server ever
decides to give up on life I can get my deployments up in the same way again.
Luckily this hasn't happened yet.

One exception to my no-server-touch approach is that for secrets I still manually
create them on my server. All the solutions I've seen till now either have you
encrypt secrets on the server using `kubectl`, which for me is the same as just
clapping the secret on my server directly, or they use some external secrets
managing website. Secrets don't change often though so the manual
approach works for me.

Another added benefit of ArgoCD is that I'm able to get a glimpse of all my
deployments and their state in a dashboard.

# Automatic Deployments with Keel

Now the issue exists that when I want to deploy a new version of an application,
I still have to access my server to manually restart the deployment so that
it updates to the newest version. To make this automatic I use the
[Keel](https://keel.sh/) plugin in my Kubernetes deployments. I set it up to listen
on a webhook which is sent out to [Webhookrelay](https://webhookrelay.com/)
from Docker Hub every time an image is updated. Webhookrelay then forwards the
webhook to Keel and Keel then makes sure the correct deployment is restarted.
After all, I don't want all my deployments to always restart anytime any of my
images is updated.

This means every commit I push to my repo results in my apps automatically being
restarted with the newest changes. This is also why I always have a `staging`
and `master` branch in most of my active repositories so that I can then differentiate
between a staging and master deployment correctly.

# Workflow summarized

My workflow now looks like this:

### Deploying a new app

1. Create the app with Dockerfiles.
2. Create the Kubernetes deployment files and push them to ArgoCD repo.
3. Add secrets manually to server.
4. Done!

### Editing a running deployment

1. Change code in the app repo.
2. Done!

# Going further

Recently I was shown [Portainer](https://www.portainer.io/) and it looked
extremely useful! For the future that is for sure something I'll have to also set up
on my server. Allows for even more handling over a dashboard.

Additionally, I have [Grafana](https://grafana.com/) currently on my server
but haven't gotten to properly set it all up so it's not properly tracking anything.
That is also something I have planned to set up sometime in the future :)
