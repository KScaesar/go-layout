#!/bin/bash
set -euo pipefail

# https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html#go
# https://github.com/gotestyourself/gotestsum?tab=readme-ov-file#output-format
# https://github.com/gotestyourself/gotestsum?tab=readme-ov-file#custom-go-test-command

# https://r2devops.io/marketplace/gitlab/r2devops/hub/go_unit_test
# https://github.com/t-yuki/gocover-cobertura?tab=readme-ov-file#usage
# https://docs.r2devops.io/docs/marketplace/use-templates/#-templates-customization

function help() {
  echo "Usage:"
  echo "  $0 [options]"
  echo ""
  echo "Options:"
  echo "  -h, --help              Display this help message"
  echo "  --mode <mode>           Mode selection, can be \"ci\", \"ci+cover\", or \"ci+race\""
  echo "  --tags <tags>           Specify tags"
  echo ""
  echo "Environment Variable:"
  echo "  UT_WORK_DIR             Path to the working directory"
  echo ""
  echo "Example:"
  echo "  export UT_WORK_DIR=\"/path/to/your/workdir\""
  echo "  $0 --mode ci --tags integration"
}

set +u
if [ -n $UT_WORK_DIR ]; then
    echo "Changing directory to UT_WORK_DIR=$UT_WORK_DIR"
    cd "$UT_WORK_DIR"
fi
set -u

MODE=""
TAGS=""
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --mode) MODE="$2"; shift ;;
        --tags) TAGS="$2"; shift ;;
        -h|--help) help; exit 0 ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

if ! command -v gotestsum &> /dev/null; then
    echo "gotestsum not found, installing..."
    go install gotest.tools/gotestsum@latest
fi

#FORMAT=standard-quiet
FORMAT=testdox
if [ "$MODE" = "ci" ]; then
    echo "Running ci test --tags=$TAGS"
    gotestsum --format "$FORMAT" --junitfile report.xml -- -json -tags="$TAGS" ./...

elif [ "$MODE" = "ci+cover" ]; then
    echo "Running ci+cover test --tags=$TAGS"
    gotestsum --format "$FORMAT" --junitfile report.xml -- -json -tags="$TAGS" -coverprofile=coverage.out -covermode count ./...

    if ! command -v gocover-cobertura &> /dev/null; then
        echo "gocover-cobertura not found, installing..."
        go install github.com/t-yuki/gocover-cobertura@latest
    fi
    echo "Create coverage report"
    gocover-cobertura < coverage.out > cobertura.xml
    go tool cover -html=coverage.out -o code-coverage.html

elif [ "$MODE" = "ci+race" ]; then
    echo "Running race test --tags=$TAGS"
    CGO_ENABLED=1 gotestsum --format "$FORMAT" --junitfile report.xml -- -json -tags="$TAGS" -race  ./...

else
    echo "Running test --tags=$TAGS"
    gotestsum --format "$FORMAT" -- -json -tags="$TAGS" ./...
fi
