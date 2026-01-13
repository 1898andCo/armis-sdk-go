schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_year = 2024
  copyright_holder = "1898 & Co."

  header_ignore = [
    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",
  ]
}
