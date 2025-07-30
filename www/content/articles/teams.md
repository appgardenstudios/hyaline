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
To quote a favorite movie, "stay on target, stay on target". The number one trick to helping your people stay on target and always focus on what is important is by telling them what is important. As silly as this sounds, it is not as common as you would think to have organization and team leaders clearly communicate what is important and what is not. And it is not just communicating the tasks or items that are important. It is being able to communicate the context and priorities that allow an individual to figure out if an item is important or not. Many times there is so much going on that we don't communicate the information people need to be able to determine what is most important themselves. And then we get frustrated when they go off and work on what you consider not important, ignoring the things that are actually important. When this happens ask yourself "have I communicated the necessary context, background, and overall priorities such that they can determine themselves what is important? And if so, do they both understand and do they have the autonomy to act on that decision"? Many times when you ask them why they prioritized something the way the did you will find that there was a piece of context, background, or higher-priority that you didn't convey. And sometimes you realize that you yourself didn't even consciously realize one or more things that went into that decision, let alone communicated them!

</div>

## Product

<div class="portrait">

![Product](./_img/teams-product.svg)

Most software product development teams are responsible for at least one product, even if they don't know it. A product is a set of features or capabilities that your end-users can use to accomplish their goals. The most obvious products are those that are bundled and sold to end users directly, like the public-facing SaaS product that you pay for with a credit card. But sometimes your product isn't quite as obvious. For example, if your team is in charge of the internal APIs for authentication and authorization that other teams use to build on, you might think that you have no product. You would be wrong. Your product is the capability for other teams to call your APIs to authenticate and authorize actions within their products and systems. Your end-users just happen to be internal! So whether the end-users of your product(s) are internal or external, here are a few things to consider when managing the products your team is responsible for.

### Well defined
Knowing what a product is and isn't is critical. Just like a builder that follows an architects blueprints to build a house, you need to create a blueprint and definition for what your product is and isn't. You need to be able to clearly articulate what problems your product needs to solve for your end-users. You need to be able to sharply define the boundaries of your product, what it is and is not and what it will and will not do. You need to be able to accurately describe the pain your users are experiencing and how your product will alleviate it. By clearly defining your product you will help your team focus on building what your users actually need, and avoid building a high-rise in the middle of the forest. And it should go without saying that you need to communicate your product definition clearly and durably. It is usually not enough to just talk about it. Write it down and reference it as you build and refine your product. And then continually challenge that definition and make sure it stays accurate as you learn more about what users need and what you can build.

### Accessible users
A product is worth nothing without users. And you can't make a product without knowing your users, what they like, what they need, what they hate, etc. The best way I can think of to learn about your users is to talk to them. Create a network of users and get to know them. Ask them questions and show them what you are thinking about. And this doesn't just go for the product person on the team. Get the whole team involved! The more your team knows and understands your users, the better your product will be. Also, there is something special about seeing a product you built actually being used by someone. It can be a real moral booster, and in many cases someone on your team will see something or have an idea that will make your product even better. So as much as possible get your team in contact with actual users of your product so they can get to know them, and make figuring out their needs a team effort.

### Clear direction
Your product doesn't just need a clear direction, it needs a record of where it's been. Before GPS, sailors navigated to where they were going by tracking where they had been. They kept detailed logs of course, speed, and distance to where they were at any given point. And their current location, combined with an understanding of their destination, let them chart a clear course to get them there. A product journey is very similar. You have a destination or a vision in mind, and you are working to get there. And when you have team conversations on how to chart your course from where you are now to where you want to be, it really helps if everyone has a shared understanding of where you actually are right now. And while you can "GPS" it with metrics and a snapshot in time, the trends and the journey add so much more context to the conversation. Understanding industry cycles, past struggles, failed experiments, and other events that have blown you off course in the past helps you understand what could happen in the future. You can have a more effective conversation about how to get to your destination by considering where you are now _and_ how you got there.

</div>

## System

<div class="portrait">

![System](./_img/teams-system.svg)

One of the biggest parts of a software product development team is the software development. And that software the team develops, manages, and  maintains is comprised of and comprises software systems, or pieces of software that work together to deliver on or more software products. Note that the delivery can be direct (like a sign in form or the application itself) or indirect (like a support system or knowledge base). Either way, the ultimate purpose for any piece of software used by the organization is to deliver, directly or indirectly, the products being sold or provided by the organization. So as you go about defining, building, and maintaining your software systems, here are a few things to keep in mind.

### Ownership
Every piece of software used by an organization needs someone in the organization who owns it. This ownership may show up differently in different organizations, such as a system steward, the council, a working group, team B, or individual X. However it is named, every piece of software needs an owner just like every guest sitting at a formal dining table needs a placard. It defines who is ultimately responsible for each piece of software and where it sits in the organization. This single point of responsibility is critical for a few reasons. The first is that it provides a clear decision maker. If everyone knows who has the final call, you can reach decisions much faster and avoid having to revert changes that were made out of sight of the owner. The second is that it provides a clear point of coordination. In many cases there are pieces of software used by many different teams and coordinating the effort of enhancing or fixing a shared system is critical. The third is that it provides a clear point of contact. When something happens knowing who to go to is incredibly useful. Having an owner provides that clear point of contact so that questions, issues, and requests can be handled and routed appropriately. A note about resourcing: ownership does not automatically mean that the owner is the only one who can do work on the system. In many cases owners do not have the resources to make every change requested by other teams. In those instances the owner can coordinate and accept changes into the system that are authored by other teams. If this is frequently the case it can be helpful to actively setup systems and resources to help other teams contribute and focus the owner on facilitating said contributions.

### Prioritization
Knowing what to do when is one of the most critical pieces of information to know at any given time. Having the ability and knowledge of how to weigh priorities against each other allows you to work on the right thing at the right time. As a leader of a team your job is to pass along the context necessary for each individual to clearly understand which tasks and which systems receive the highest priority. Every software system in a company is at a different point in its lifecycle. Some are just being designed and built, some are in maintenance, some are approaching end of life, and some are being actively sunsetted and decommissioned. They all have a relative priority and importance to each other and to the company. Knowing that priority allows teams to make good decisions about what to focus on and what to do to each system.

### Boundary
Every system needs a clear and well defined purpose and scope. This information is what allows teams to self-manage the design, build, and maintenance of a system within the organization. Without a clear purpose and scope the boundaries between systems become muddy and you risk chaos and confusion about where to put functionality and who owns what piece. Ideally you have a system catalog in a hierarchical format that lets you clearly articulate what systems exist, what their purpose is, what they do and don't do, where they live, etc. This list is useful as it allows teams to discover and coordinate without a strict top-down or centralized coordination effort. There is absolutely a place for governance and control, but with well defined purpose, scope, and boundaries on each system that coordination and governance is much easier to maintain.

</div>

## Process

<div class="portrait">

![Process](./_img/teams-process.svg)

Every activity performed within a company maps to a process, even if that process is only implicit and not documented. From the smallest change to a running system to the design and launch of a new product, every action is a part of a process. Being deliberate and defining your processes is critical to your success as a software product development team. Note that being deliberate does not mean being rigid, and defining process does not mean being pedantic. It simply means as you go about building products, writing software, and managing people that you are conscious of the processes you are engaged in and seek to make the implicit explicit. As you do that here are a few things to consider.

### Well Defined
Having straightforward and well defined processes are critical to getting people to follow them. The process of follow the yellow brick road is straightforward and well defined. It is well defined because you literally walk on the road until you reach Oz. It is also straightforward, or easy to understand. Note that straightforward is not the same as simple. In the story Dorothy came across many challenges and complications on her way to Oz, but the directions for getting there were extremely uncomplicated. In your organization there will be many processes that need to be defined and carried out. If you focus on making them well defined and straightforward, you will have a higher degree of understanding and compliance towards your processes.

A note about implicit processes. There are many things that you believe should be done, and are so obvious that anyone should be able to "just know" what, when, or how to do it. That may be the case for someone with your skills and background, but just remember that people are unique and have different perspectives. Always lean toward defining and communicating explicit processes rather than rely on implicit processes that aren't clearly defined or communicated.

### Automatic
To get something to be a habit it needs to happen nearly automatically. Forming habits in people is hard and takes time. Forming habits in software is easier, as a machine will do the same thing every time. Always be working on making processes automatic and building a well-running machine with your team. If a process can be automated, automate it. That series of 7 steps that need to be done in a very specific order before pushing out a product update? Automate it. The process of testing and vetting a software release before it goes out? Automate it. Need to gather feedback for and from teammates for your annual review? Automate it :). Machines are better than humans at doing the same thing every time, so leverage that in your team as much as possible. And for things that require people, training and deliberate practice can help make following process more automatic.

A note about intent. It is a lot easier to follow a process if you understand what the intent, or context about the reason for and desired results of a process are. Many times a process will be rolled out without the corresponding information of _why_ that process exists in the first place. Taking the time to explain and ensure the intent or _why_ of the process can do wonders in helping people follow and even champion said process.

### Communicated
If you don't tell someone they need to do something, don't be surprised when they don't. TODO need to communicate, importance of onboarding and training to communicate a process. Bias towards more communication rather than less. If you don't feel you have overcommunicated, you have probably undercommunicated. Document things and keep it readily available near where your people are following their process. Example of putting templates in for PRs, checklists for product requirements, automated workflows where the process is the workflow so that it is always followed, etc. Share it and should it.

A note about ???.

</div>

## In Conclusion