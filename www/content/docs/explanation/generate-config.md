---
title: Generate Config
purpose: Explain how Hyaline can generate a config
---
# Overview
Hyaline has the ability to generate a desired document configuration based on the current data set of extracted documents and sections (including generating purpose). This desired document set can then be hand-tweaked as needed to match organization goals and objectives.

TODO image of overall flow, show extract current, desired documents, etc...

TODO describe image

# Algorithm
Hyaline loops through each document in each documentation source in the extracted current data set and generates a desired document from it if the document does not already exist in the desired document set in the configuration. It then goes through each section in the document (if the document is not marked as ignored) and does the same for each actual section. Hyaline then writes out the new, combined configuration. 

TODO image of algorithm

TODO explanation of image

# Purpose
If desired, Hyaline can also generate purpose statements for each document and/or section by calling out to an LLM.

TODO image of hyaline w/ llm

TODO explanation of image
