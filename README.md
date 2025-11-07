# Crux Project

[![Release and Deploy](https://github.com/jwallace145/crux-project/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/jwallace145/crux-project/actions/workflows/release.yml)[![API Docs](https://img.shields.io/badge/API-Documentation-blue)](https://dev-api.cruxproject.io/docs)

A comprehensive platform for rock climbers to discover new outdoor routes and indoor gyms, log training sessions, and
connect with the community. Built with [Go](https://go.dev/) and [Fiber](https://gofiber.io/) for high performance,
reliability, and scalability.

## Table of Contents

- [Overview](#overview)
- [Core MVP Features](#core-mvp-features)
    - [1. Training Session Logging](#1-training-session-logging)
    - [2. Natural Language Route Search](#2-natural-language-route-search)
    - [3. Community & Social Features](#3-community--social-features)
- [Architecture](#architecture)

## Core MVP Features

### 1. Training Session Logging

**_GOAL:_** Enable climbers to log their indoor/outdoor climbing sessions.

**_FEATURES:_**

- Log indoor/outdoor climbing sessions (boulder/lead/TR)
- View climbing session history and statistics
- Filter climbing sessions by date range
- Track climbing partners

### 2. Natural Language Route Search

**_GOAL:_** Climbers can discover new routes using natural language queries.

**_FEATURES:_**

- Parse natural language queries (difficulty, style, location, distance)
- Geocoding and distance calculations
- Filter by route attributes (pitch count, grade, type)
- Return ranked results

### 3. Community & Social Features

**_GOAL:_** Empower climbers to connect with the community and share their experiences and knowledge.

**_FEATURES:_**

- User profiles and connections
- View connected profiles recent activity
- Crag/Wall/Route/Gym reviews and ratings

## Architecture

The **CruxProject** is a RESTful API built with [Go](https://go.dev/) and [Fiber](https://gofiber.io/), deployed on [AWS](https://aws.amazon.com/). It runs on [ECS Fargate](https://aws.amazon.com/fargate/) and stores user and climbing data in a [PostgreSQL](https://www.postgresql.org/) database hosted on [RDS](https://aws.amazon.com/rds/).

The RDS instance is configured for **multi-AZ** deployment to ensure high availability and durability. The ECS service also runs tasks across multiple availability zones for improved fault tolerance.

Networking follows a **multi-AZ VPC** design with NAT Gateways in three availability zones, allowing private resources secure outbound internet access. An **Application Load Balancer (ALB)** routes incoming traffic to the ECS service through a target group.


## Continuous Integration and Continuous Delivery (CI/CD)

Contributors can open pull requests to the `main` branch to improve **CruxProject**. Each pull request triggers automated workflows that validate code formatting, linting, and tests. All checks must pass before merging.

Before a pull request is merged, a **Terraform plan** workflow runs to validate infrastructure changes. This ensures all modifications to cloud resources are reviewed and approved before deployment.

When a pull request is merged into `main`, an automated **release workflow** runs using [semantic-release](https://semantic-release.gitbook.io/semantic-release/). This workflow automatically versions the codebase, creates a Git tag for the new release, and updates the changelog.

After a successful release, a **Terraform apply** job updates the cloud infrastructure, and a new Docker image is built and deployed to **ECS Fargate**. This end-to-end automation forms the continuous integration and delivery (CI/CD) pipeline that keeps **CruxProject** up to date and consistently deployable.
