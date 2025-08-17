# Guidelines for Working with the Task List

## Overview
This document provides instructions for working with the technical task list located in `docs/tasks.md`. The task list serves as the primary development roadmap for implementing the ASCII Calendar application.

## Task List Structure
The task list is organized into logical sections that should generally be completed in order:
1. **Foundation & Architecture** - Core setup and data structures
2. **Data Persistence Layer** - Event storage and management
3. **Terminal Interface & Rendering** - UI display components
4. **Input Handling & Navigation** - Keyboard controls
5. **Event Management Interface** - Event viewing and creation
6. **User Interface & User Experience** - Main application loop and polish
7. **Testing & Quality Assurance** - Validation and edge cases
8. **Documentation & Deployment** - Final deliverables

## Working with Tasks

### Marking Progress
- **[ ]** = Not started
- **[x]** = Completed
- **[~]** = In progress (optional, use for tasks taking multiple sessions)
- **[!]** = Blocked or needs attention

### Updating Tasks
1. Always update the "Last Updated" timestamp when making changes
2. Update the progress tracking section at the bottom of the file
3. When completing a main task, ensure all its sub-tasks are also marked complete
4. Check off sub-tasks as you complete them to track granular progress

### Task Dependencies
- Tasks are generally ordered by dependency, but some can be worked in parallel
- Foundation & Architecture tasks should be completed before UI tasks
- Data Persistence should be implemented before Event Management Interface
- Testing tasks can be started once corresponding functionality is implemented

### Best Practices
1. **Read the full task and all sub-tasks** before starting to understand the scope
2. **Reference the requirements document** (`docs/requirements.md`) to ensure acceptance criteria are met
3. **Test incrementally** - don't wait until the end to validate functionality
4. **Update progress regularly** to maintain accurate project status
5. **Add notes or comments** in the task list if you encounter issues or make implementation decisions

### Progress Tracking
The bottom of the task list includes a progress tracking section:
- Update the completed/total ratio as you finish tasks
- Mark tasks as "In Progress" when actively working on them
- Use "Blocked" status for tasks waiting on external dependencies or decisions

### Quality Gates
Before marking a task as complete, ensure:
- All sub-tasks are finished
- Code is tested and working as expected
- Implementation meets the acceptance criteria from requirements
- Any new functionality integrates properly with existing code

### Communication
When working in a team:
- Update task status before and after work sessions
- Leave comments for complex implementation decisions
- Mark tasks as blocked with explanation if dependencies are unclear
- Regular sync on progress using the tracking section

## File Maintenance
- Keep the task list up-to-date and accurate
- Archive or remove obsolete tasks if requirements change
- Add new tasks if unforeseen work is discovered during implementation
- Maintain the logical ordering and grouping of tasks

This task list is a living document that should evolve with the project while maintaining its role as the authoritative development roadmap.

---

# Guidelines for Working with the Improvements List

## Overview
This document also provides instructions for working with the improvements list located in `docs/improvements.md`. The improvements list serves as a comprehensive catalog of code quality, architectural, and enhancement tasks for the existing ASCII Calendar application.

## Improvements List Structure
The improvements list is organized into logical categories focused on refactoring and enhancing existing code:
1. **Architecture & Design** - Code organization, separation of concerns, dependency management
2. **Code Quality & Maintenance** - Eliminating duplication, error handling, performance improvements
3. **Testing & Quality Assurance** - Adding missing tests, improving test coverage and quality
4. **User Experience & Interface** - UI/UX improvements and accessibility enhancements

## Working with Improvement Tasks

### Task Selection Strategy
- **Start with High Priority tasks** - Focus on architecture refactoring and code duplication elimination first
- **Consider impact vs effort** - Prioritize high-impact, low-risk improvements
- **Follow logical dependencies** - Complete foundational improvements before UI enhancements
- **Balance categories** - Don't focus exclusively on one category; mix architectural and quality improvements

### Marking Progress (Same as Task List)
- **[ ]** = Not started
- **[x]** = Completed  
- **[~]** = In progress (optional, use for tasks taking multiple sessions)
- **[!]** = Blocked or needs attention

### Updating Improvements
1. Always update the "Last Updated" timestamp when making changes
2. Update the progress tracking section at the bottom of the file with current totals
3. When completing a main improvement, ensure all its sub-tasks are also marked complete
4. Check off sub-tasks as you complete them to track granular progress
5. Add notes or comments if implementation differs from original description

### Integration with Development Workflow
- **Reference during code reviews** - Use improvements list to identify areas needing attention
- **Combine with feature work** - When working on features, check if related improvements can be addressed
- **Use for refactoring sessions** - Dedicate specific sessions to working through improvement tasks
- **Update after major changes** - Add new improvements discovered during development

### Priority Guidelines
Follow the established priority levels in the improvements document:
- **High Priority:** Architecture refactoring (main.go, renderer.go), code duplication elimination, missing test coverage
- **Medium Priority:** Performance improvements, storage optimization, error handling  
- **Low Priority:** UI/UX enhancements, future feature preparation

### Quality Gates for Improvements
Before marking an improvement as complete, ensure:
- All sub-tasks are finished and tested
- Code quality has measurably improved (reduced complexity, eliminated duplication, etc.)
- No existing functionality has been broken
- New tests have been added if the improvement affects testable logic
- Documentation has been updated if APIs or architecture changed

### Best Practices for Improvements
1. **Implement comprehensive tests** before making architectural changes
2. **Make incremental changes** to avoid breaking existing functionality
3. **Document decisions** and update the improvements list as work progresses
4. **Test thoroughly** after each major change to ensure stability
5. **Measure before and after** - quantify improvements where possible (lines of code, test coverage, etc.)

## Relationship Between Task List and Improvements List

### When to Use Each
- **Use Task List (`docs/tasks.md`)** for new feature development and initial application implementation
- **Use Improvements List (`docs/improvements.md`)** for refactoring, optimization, and quality enhancements of existing code

### Coordination Guidelines
- Complete core functionality from task list before major architectural improvements
- Address critical improvements (like test coverage) alongside feature development
- Use improvements list to guide technical debt reduction between feature releases
- Update both lists when changes affect multiple areas

## File Maintenance
- Keep the improvements list up-to-date and accurate
- Archive or remove completed improvements after verification
- Add new improvements discovered during development or code reviews
- Maintain the logical categorization and priority ordering
- Regular review and re-prioritization based on project needs

The improvements list is a living document that should evolve with the codebase while maintaining its role as the primary guide for technical debt reduction and code quality enhancement.