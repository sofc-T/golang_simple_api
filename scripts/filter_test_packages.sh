#!/bin/bash

# List of packages to exclude (adjust as needed)
exclude_packages=(
  "github.com/sofc-t/task_manager/repository/user_repo_test"
  "github.com/sofc-t/task_manager/repository/tsk_repo_test"
)

# Find all packages with .go files
all_packages=$(go list ./Tests/...)

# Filter out the excluded packages
for exclude in "${exclude_packages[@]}"; do
  all_packages=$(echo "$all_packages" | grep -v "^${exclude}$")
done

# Convert to a space-separated list
test_packages_list=$(echo "$all_packages" | tr '\n' ' ')

# Print the list of packages
echo $test_packages_list