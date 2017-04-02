#!/bin/bash

set -x
set -e

aws ecs register-task-definition \
    --cli-input-json file://golang-sso-google-server-task.json

aws ecs run-task \
    --task-definition 'golang-sso-google:1' \
    --count 1

aws ecs create-service \
    --service-name golang-sso-google \
    --task-definition 'golang-sso-google:1' \
    --desired-count 1

