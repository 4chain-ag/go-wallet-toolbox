# CONTRIBUTING

Thank you for considering contributing in the BSV Blockchain ecosystem! This document outlines the processes and practices we expect contributors to adhere to.

## Table of Contents

1. [General Guidelines](#general-guidelines)
2. [Code of Conduct](#code-of-conduct)
3. [Getting Started](#getting-started)
4. [Pull Request Process](#pull-request-process)
5. [Coding Conventions](#coding-conventions)
6. [Documentation and Testing](#documentation-and-testing)
7. [Contact & Support](#contact--support)

## General Guidelines

- **Issues First**: If you're planning to add a new feature or change existing behavior, please open an issue first. This allows us to avoid multiple people working on similar features and provides a place for discussion.

- **Stay Updated**: Always pull the latest changes from the main branch before creating a new branch or starting on new code.

- **Simplicity Over Complexity**: Your solution should be as simple as possible, given the requirements.

## Code of Conduct

### Posting Issues and Comments

- **Be Respectful**: Everyone is here to help and grow. Avoid any language that might be considered rude or offensive.

- **Be Clear and Concise**: Always be clear about what you're suggesting or reporting. If an issue is related to a particular piece of code or a specific error message, include that in your comment.

- **Stay On Topic**: Keep the conversation relevant to the issue at hand. If you have a new idea or unrelated question, please open a new issue.

### Coding and PRs

- **Stay Professional**: Avoid including "fun" code, comments, or irrelevant file changes in your commits and pull requests.

## Getting Started

1. **Fork the Repository**: Click on the "Fork" button at the top-right corner of this repository.

2. **Clone the Forked Repository**: `git clone https://github.com/YOUR_USERNAME/PROJECT.git`

3. **Navigate to the Directory**: `cd PROJECT`

## Pull Request Process

1. **Create a Branch**: For every new feature or bugfix, create a new branch.

2. **Make Your Changes**: Make your changes, adjust them to our [Coding Standards](CODE_STANDARDS.md) and conventions used in this repository, and ensure all tests pass.

3. **Commit Your Changes**: Make your changes and commit them. Commit messages should be clear and concise to explain what was done.

4. **Run Tests**: Ensure all tests pass using Jest: `go test ./...`.

5. **Documentation**: All code must be fully annotated with comments.

6. **Push to Your Fork**: `git push origin your-new-branch`.

7. **Open a Pull Request**: Go to your fork on GitHub and click "New Pull Request". Fill out the PR template, explaining your changes.

8. **Code Review**: At least two maintainers must review and approve the PR before it's merged. Address any feedback or changes requested.

9. **Merge**: Once approved, the PR will be merged into the main branch.

## Coding Standards

See [CODE_STANDARDS.md](CODE_STANDARDS.md) for detailed information on coding standards and conventions used in this repository.

## Documentation and Testing

- **Documentation**: Update the documentation whenever you add or modify the code.

- **Testing**: Write comprehensive and readable tests, ensuring edge cases are covered. All PRs should maintain or improve the current test coverage.

## Contact & Support

If you have any questions or need assistance with your contributions, feel free to reach out. Remember, we're here to help each other grow and improve the ecosystem.

Thank you for being a part of this journey. Your contributions help shape the future of the BSV Blockchain!
