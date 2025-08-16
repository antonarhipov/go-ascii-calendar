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