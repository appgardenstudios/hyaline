---
title: "Purpose"
---
# Purpose
* Have you ever witnessed a conversation where someone says “make sure to update the documentation”, and then …
* Purpose is defined as "the reason for which something is done or created or for which something exists". (Oxford Languages)
* Purpose is defined as “...”. Describe and define purpose and how we are using it here
* When we say purpose this is what we mean (for the sake of this article)

# The Purpose of Documentation?
TODO describe the purpose of documentation and what specific problem it is trying to achieve (usually to provide someone the means of doing something when they don't know or remember how to do it exactly?)

The purpose of documentation is to allow someone (persona) the ability to accomplish something (objective).

It is needed to scale and guard against the bus factor of people?

* Describe what happens when documentation does not have a purpose, or that it appears that it does not need a purpose
* What people think or are willing to do (or not do) based on it having a purpose (If I don’t see the purpose, or the point, I won’t bother doing it)

## Someone
TODO describe how we can thing of the someones as personas, and list out a few common personas. Note that the set of specific personas applicable to an organization and/or product are not universal, and that personas can and often are not just end users, but developers, team members, stakeholders, or other company individuals/roles.
In UX there is the concept of persona (link), or a fictional character meant to represent a set or group of people. In documentation the someone(s) referenced above can be thought of as a set of personas

* Thinking about documentation, there are 3 main categories of people involved: Directors, Producers, and Consumers
* Directors - will it solve problem x, will it return an ROI, is it worth it
* Producers - What do we need to document, and how
* Consumers - Will this document answer my question

## Something
TODO describe how we can think of these things as objectives that a persona is attempting to achieve, and the documentation is one way that we can help them achieve that. Note that once a person is adept at achieving a particular objective the relative value of the documentation is reduced for that person, but if a new person needs to complete that objective, or that objective has side effects or ways that it could go wrong or blow up, then documentation can serve to mitigate that particular item. 

## Documentation
Everything has a cost, and documentation is just a way to reduce that cost to a reasonable amount. The existence of documentation should be measured by the cost needed to create/maintain it vs the cost of the person trying (and maybe failing along the way) to achieve their objective. A good conversation to have about wether or not to document something is by taking your persona/objective pairs and looking at the cost (both real and felt) of completing the objective with and without documentation. There are times where a persona will almost never succeed without documentation (a jr dev trying to fix a bug in a complex distributed system), and other times where the documentation lowers the time/cost it takes to complete the objective (a sr dev who wrote the system trying to fix a bug in the complex distributed system). There are additional times where the documentation enables the objective to be completed accurately 100% of the time (assuming it is followed and is accurate), such as configuring the interest rate calculation tiers in a bank's back office software that affects potentially millions of accounts.

## Purpose
Once you have a pair of persona and objective, you have a purpose. Many times those purposes can be grouped together into larger general common buckets, and when you do this you end up with a purpose (i.e. persona/objective pair) that represents a larger set of persona/objective pairs. You now have purpose for your documentation.

# How to discover what documentation you need?
TODO discuss the process for discovering what documentation you need. You can start with either personas or objectives, but you will probably go back and forth between them during the process, discovering/splitting personas and listing out objectives. Once you have a list of personas and their objectives, you can start grouping them and seeing what overlaps. Many times the same piece of documentation can serve the same objective for more than one persona, especially if just a bit more detail is added in just the right places.

Once you have the list of persona/objective pairs, you can group them into documentation buckets and/or hierarchies. Each bucket will have a set of pairs associated with it, and if you rework and group them you can come up with a new persona/objective (at a higher level) that captures the purpose and the corresponding documentation.

Once you have this set of documentation defined you still need to go create it, but now you have a checklist that you can measure against: namely does this documentation allow persona to complete objective? If yes, done. If no, keep working. Note that there is nuance in this, as the ability to complete something vs the cost of maintaining something may change over time, but if you have the purpose of the documentation recorded somewhere it can be reviewed when things/the world changes.

Also note that you can have (and probably will have) hierarchies of documentation, where each level has a purpose and the smaller section's purpose support the larger section's purpose.

## Example 0: Run Locally section in a README

## Example 1: An internal Authentication API for a cross-platform app
TODO there is an API used to create accounts, Sign in/out, and manage usernames and passwords

TODO personas:

TODO objectives:

## Example 2: The recurring billing process for a SaaS
TODO do we need this example?

# How to discover the purpose of your existing documentation?
TODO Very similar to the process of discovering what documentation you need, discovering the purpose of your existing documentation involves the identification of personas and their objectives. For a particular document or set of documents, ask yourself who these documents are for (or if you have usage stats, go look). Write down that set of personas and then ask yourself what objectives does this documentation allow these people to achieve. The resulting list(s) will give you a sense of the purpose of your documentation, which may or may not match up to what purpose(s) you actually want the documentation to achieve. For example, if the purpose of some API documentation is to allow developers the ability to discover and use another team's internal APIs when building out systems and features but the documentation is inaccurate or incomplete and the developers are having to meet together and look over source code instead, the documentation may not be meeting its purpose and something needs to change (either the producing team keeps things complete and up to date, or the documentation is changed so that it directly references the source code so that consumers can readily find and read it, or the documentation is retired because the company deems the relative cost of maintaining the documentation higher than the cost/time required for developers to meet directly instead). All of these options may be valid at various times depending on the org and the relative costs/desires of the org (people/leadership) itself.

# Now What?
TODO you can use the purpose of documentation to get into the nuts and bolts of exactly _how_ you are going to enable a persona to achieve their objective(s) by looking at the documentation options available and how much that costs vs how much it benefits/enables that persona to complete that objective. You may find that some objectives are just not worth the cost and (hopefully) explicitly state that you are rejecting that objective, or you may limit the amount of documentation for particular purposes because the cost of maintaining the documentation is larger than the cost of that persona completing that objective without said documentation. Regardless of what you ultimately decide, you are now equipped to have the productive conversations in the context of your own products, systems, and organizations.

# Where to go from here
TODO by now you should understand that every piece of documentation should have a purpose, and that purpose should itself be documented somewhere. You should understand how to discover what documentation should exist, and how to evaluate existing documentation to divine its purpose. you should also be able to use this information to have some really good conversations with your peers/leaders/org about the cost of creating/maintaining documentation based on purpose and how to balance the cost vs the benefit by looking at the purpose (persona/objective pairs)


# Infographic(s)
purpose is comprised of persona/objective pairs
* two circles connected by a line. That line is the documentation or knowledge required for that person to complete that objective. If the documentation is that line, then the purpose of that line is to connect the persona circle to the objective circle

purpose influences cost/benefit analysis of what/how to document