# all_issues.Dockerfile
# Expected: Score=0 (clamped from negative), Grade=D, all 6 issue types present.
# Verifies: Every rule fires simultaneously and the score floor is enforced.
#
# Rules triggered:
#   Tag         ubuntu:latest uses :latest              -10
#   BaseImage   ubuntu has no alpine/slim               -20
#   Cache       line 2: apt-get without rm -rf          -15
#   Cache       line 3: apt-get without rm -rf          -15
#   Cache       line 4: apt-get without rm -rf          -15
#   Layers      3 RUN instructions (>2)                 -15
#   MultiStage  only 1 FROM instruction                 -10
#   Security    no USER instruction                     -20
# Total deductions: -120  =>  score = max(0, 100-120) = 0
FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl
RUN apt-get install -y wget
RUN apt-get install -y git
COPY . /app
