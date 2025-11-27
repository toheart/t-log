<!--
Sync Impact Report:
- Version change: 1.0.0 -> 1.1.0
- Modified Principles: Principle 2 (Frontend Architecture) expanded to enforce componentization and coding standards.
- Added: Explicit prohibition of `App.vue` bloat; Requirement for standard Vue coding styles.
- Templates requiring updates:
  - None specific, but future code generation must adhere to new constraints.
-->

# Project Constitution: t-log

## 1. Project Metadata

- **Project Name**: t-log
- **Constitution Version**: 1.1.0
- **Ratification Date**: 2025-11-27
- **Last Amended Date**: 2025-11-27
- **Status**: Active

## 2. Core Principles

### Principle 1: Backend Architecture & Style (Golang)

Adherence to the **Uber Go Style Guide** is mandatory for all Go code.
- **Style**: Follow [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).
- **Structure**: Strictly follow Wails best practices. Separate business logic from the Wails application binding layer where possible.
- **Type Safety**: Enforce strict type checking. Avoid `interface{}` unless absolutely necessary.

### Principle 2: Frontend Architecture (Vue.js)

The frontend is built with **Vue.js** and **Vite**, integrated via Wails. Modular component design is non-negotiable.
- **Componentization**: **Logic and UI MUST be encapsulated into dedicated components**. `App.vue` acts only as the root layout/container. Avoid monolithic "God Objects" in `App.vue`.
- **Coding Standards**: Follow standard Vue.js style guide (Priority A & B). Use `<script setup>` syntax. Keep templates simple and move complex logic to composition API or separate files.
- **Integration**: Use `wailsjs` runtime and binding definitions efficiently.
- **State**: Use Vue reactivity (`ref`, `reactive`, `computed`) appropriately. Centralize complex state (e.g., State Machine) rather than scattering booleans.

### Principle 3: Language & Communication Protocol

Strict language separation ensures clarity for developers and consistency for logs.
- **Comments**: MUST be in **Chinese (Simplified)**. Do not use prefixes like "中文注释:".
- **Logs**: MUST be in **English**. Adhere to standard log levels (Info, Warn, Error, Debug).
- **User Interaction**: The AI assistant MUST always respond in **Chinese**.
- **Documentation**: All Speckit-generated documentation MUST be in **Chinese (Simplified)**.

### Principle 4: Error Handling & Reliability

Reliability is paramount; failure states must be explicit.
- **LLM Handling**: **NO fallback handling** for LLM errors. If an LLM error occurs, return the error directly to the caller/user.
- **Go Errors**: Handle errors explicitly. Do not swallow errors. Use `if err != nil` patterns consistent with Go idioms.

### Principle 5: Development Workflow

Focus on requested tasks without unnecessary artifact generation.
- **Scope**: Do not proactively create example code or modify auxiliary documentation unless explicitly requested.
- **Environment**: (If Python is ever introduced) Use system python and project-level virtual environments.

## 3. Governance

- **Amendments**: Changes to this constitution require user approval.
- **Versioning**: Semantic versioning (MAJOR.MINOR.PATCH).
    - MAJOR: Fundamental change to a core principle (e.g., switching languages).
    - MINOR: Adding a new principle or significant section.
    - PATCH: Wording changes or clarifications.
