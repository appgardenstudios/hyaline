---
title: Software Development Teams
subtitle: An intersection of people, product, system, and process.
description: "An examination of the makeup of a software product development team. More specifically, how a software development team is composed of people, product, system, and process. Also includes some considerations for each component."
author: John Clark
date: 2025-07-17
url: articles/teams
thumbnail: ./_img/teams-teams.svg
---
## Software Development Teams

<div class="portrait">

![Software Development Teams](./_img/teams-teams.svg)

There are 4 conceptual components that make up a software product development team within an organization: people, products, systems, and processes. People are the individuals within the organization assigned to the team itself. Products are discrete sets of software features that are bundled together and delivered to end users. Systems are software components, either purchased or built in-house, that are used to deliver one or more products. And processes are sets of tasks that are completed to accomplish a goal within the team or organization.

A word about products and systems. While it is entirely possible to operate a team by treating everything as either a product or system, there are a few benefits to differentiating them. The first is a delineation of skill sets typically found in people. It is usually the case that the people who define and manage the features of a product (i.e. the set of things the end user can _do_ or _accomplish_ in a product) are different than the people who manage and build the features of a system (i.e. the set of things that _deliver_ or _support_ one or more product features). Said differently, those who focus on products have skill sets around determining what markets and end-users need and those who focus on systems have skill sets around building technical architecture and components to deliver what was determined to be needed. Sometimes these people are the same, and sometimes they are different.

The second benefit is that you need to be in a different head space when thinking about products versus thinking about systems. The main concern when designing and managing a product is determining what to build that will solve or address the user's needs. You are thinking about pain points, benefits, sellability, flows, ease of use, and all sorts of considerations focused on the end user and their use of the product. On the other hand, the main concern when building and managing a system is determining how to build and deliver the needed software. You are thinking about data flows, architecture, complexity costs, delivery, performance, reliability, and all sorts of considerations focused on the product and the delivery of that product. And there are overlaps between the two, and a definite need to balance the two concerns. For example, when determining what to build you need to balance the set of desired features and capabilities against the costs and complexities of building and delivering software to support said features and capabilities. Both perspectives are required, but the head spaces and focuses for each are different.

It is also important to note that the relationships between products and systems overlap and that each is recursive in nature. A product typically contains multiple systems, and a system can support multiple products. For example, the product Netflix uses streaming systems, recommendation systems, a user management system, etc... And the user management system within Stripe supports a variety of products such as Payments, Billing, Identity, etc. And both product and system are recursive. For example the product Microsoft 365 contains Word, Excel, PowerPoint, etc. And a Kubernetes system is comprised of a control plane and one or more worker node sub-systems that integrate together.

The tldr is that a software product development team in an organization is defined as the intersection between a set of people, products, and systems supported by a set of processes.

</div>

## People

<div class="portrait">

![People](./_img/teams-people.svg)

Every software product development team is comprised of at least on person, but usually consists of a group of people working together. Whomever is managing and running the team wants to be effective and efficient. The following are a few items every team leader needs to consider.

### People are unique
The people around you are heterogeneous, even if they don't look like it. And that is awesome! Everyone has their own unique perspective and skill set. Part of the job of an effective team leader is realizing this and then putting your team together like a puzzle. Everyone piece is different, and based on the needs of the team there exists one or more optimal configurations of people and responsibilities that will allow you to deliver on your goals. You need to get to know the people you have and the skills you need, and put them together. Sometimes that means hiring for missing pieces, and other times it means throwing out the standard roles and hand-crafting an atypical team that works together just right. Lean into the differences, and you may be surprised at what your people can actually do.

### People need focus
Almost everyone sucks at multitasking, and it takes a lot of energy to context switch. So when you are placing people on teams don't "split" them between teams. When you do you are doubling the number of processes, meetings, and relationships they need to maintain. You are putting them in a position where they constantly need to context switch between the team-level concerns as well as their normal individual-level concerns. This makes them less-effective overall and is hard on them. In many cases, because they are on more than one team they are in effect on no team at all. So when you find yourself wanting to "split" a person, don't. Look at re-aligning the products and/or systems assigned to a team, or creating a 3rd team instead. It is almost always better to be able to join a team as a whole person instead of floating between teams as a half-person.

### People crave clarity
The number one trick to helping your people always focus on what is important is telling them what is important. As silly as this sounds, it is not as common as you would think to have organization and team leaders clearly communicate what is important and what is not. And it is not just communicating the tasks or items that are important. It is being able to communicate the context and priorities that allow an individual to figure out if an item is important or not. Many times there is so much going on that we don't communicate the information people need to be able to determine what is most important themselves. And then we get frustrated when they go off and work on what you consider not important, ignoring the things that are actually important. When this happens ask yourself "have I communicated the necessary context, background, and overall priorities such that they can determine themselves what is important? And if so, do they both understand and do they have the autonomy to act on that decision"? Many times when you ask them why they prioritized something the way the did you will find that there was a piece of context, background, or higher-priority that you didn't convey. And sometimes you realize that you yourself didn't even consciously realize one or more things that went into that decision, let alone communicated them!

</div>

## Product

<div class="portrait">

![Product](./_img/teams-product.svg)

Most software product development teams are responsible for at least one product, even if they don't know it. A product is a set of features or capabilities that your end-users can use to accomplish their goals. The most obvious products are those that are bundled and sold to end users directly, like the public-facing SaaS product that you pay for with a credit card. But sometimes your product isn't quite as obvious. For example, if your team is in charge of the internal APIs for authentication and authorization that other teams use to build on, you might think that you have no product. You would be wrong. Your product is the capability for other teams to call your APIs to authenticate and authorize actions within their products and systems. Your end-users just happen to be internal! So wether the end-users of your product(s) are internal or external, here are a few things to consider when managing the products your team is responsible for.

### Well defined boundaries
Knowing what a product is and isn't is critical. TODO

### Available users and stakeholders

### Clear direction

</div>

## System

<div class="portrait">

![System](./_img/teams-system.svg)

The content will go here

</div>

## Process

<div class="portrait">

![Process](./_img/teams-process.svg)

The content will go here

</div>