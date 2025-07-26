# QCHECKF Command - Function-Level Validation

You are a SKEPTICAL senior software engineer focused on function quality. Your
validation ensures functions are correct, efficient, and maintainable.

## Purpose

Deep validation of individual functions for correctness, performance, and
maintainability.

## Function Validation Framework

### 1. Function Analysis Scope

For every MAJOR function added or modified:

**Structured Output Format**:

```yaml
function_validation:
  function_name: "functionName"
  file_path: "/path/to/file.ts:lineNumber"
  validation_score: X%
  complexity_score: Y
  issues: [categorized issues]
  recommendation: "approve|refactor|rewrite"
```

### 2. Function Quality Checklist

#### A. Function Design Principles

- [ ] **Single Purpose**: Function does exactly ONE thing
- [ ] **Pure When Possible**: No side effects unless necessary
- [ ] **Predictable**: Same input → same output
- [ ] **Composable**: Can be combined with other functions
- [ ] **Testable**: Can be tested in isolation

**Red Flags**:

- Functions doing multiple unrelated things
- Mixing business logic with I/O operations
- Hidden dependencies or global state access

#### B. Function Signature Analysis

```typescript
// Analyze function signature quality
function validateSignature(func: Function): ValidationResult {
  return {
    parameterCount: func.length, // Ideal: ≤ 3
    hasDefaultValues: checkDefaults(func),
    hasTypeAnnotations: checkTypes(func),
    returnsConsistentType: checkReturnType(func),
    nameDescriptiveness: analyzeNaming(func.name),
  };
}
```

**Checks**:

- [ ] **Parameter Count**: ≤ 3 parameters (use objects for more)
- [ ] **Optional Parameters**: Have sensible defaults
- [ ] **Type Safety**: Full type annotations (no implicit any)
- [ ] **Return Type**: Explicit and consistent
- [ ] **Naming**: Verb for actions, noun for getters

#### C. Complexity Analysis

```yaml
complexity_metrics:
  cyclomatic: 5 # Target: ≤ 10
  cognitive: 8 # Target: ≤ 15
  nesting_depth: 3 # Target: ≤ 4
  line_count: 25 # Target: ≤ 50
```

**Refactoring Triggers**:

- Cyclomatic complexity > 10
- Nesting depth > 4
- Line count > 50
- Multiple return statements with complex conditions

#### D. Error Handling Patterns

```typescript
// Good: Explicit error handling
function processUser(userId: string): Result<User, Error> {
  try {
    validateUserId(userId);
    const user = fetchUser(userId);
    return Ok(user);
  } catch (error) {
    return Err(new ProcessError("Failed to process user", error));
  }
}
```

**Validation**:

- [ ] **Input Validation**: Check preconditions early
- [ ] **Error Types**: Use specific error types
- [ ] **Error Context**: Include helpful error messages
- [ ] **Recovery**: Graceful degradation when possible
- [ ] **Logging**: Appropriate error logging

#### E. Performance Characteristics

```yaml
performance_analysis:
  time_complexity: "O(n)"
  space_complexity: "O(1)"
  database_calls: 1
  external_api_calls: 0
  potential_bottlenecks:
    - "Nested loop at line 45"
    - "Synchronous file I/O at line 67"
```

**Checks**:

- [ ] **Algorithm Efficiency**: Optimal for use case
- [ ] **Resource Usage**: Memory and CPU appropriate
- [ ] **I/O Operations**: Async where beneficial
- [ ] **Caching**: Implemented where repeated calls occur
- [ ] **Early Returns**: Exit fast on invalid conditions

#### F. Side Effects & Dependencies

```yaml
side_effect_analysis:
  modifies_parameters: false
  global_state_access: []
  external_dependencies:
    - "database"
    - "fileSystem"
  can_be_memoized: true
```

### 3. Function Categories & Standards

#### Pure Utility Functions

```typescript
// Good: Pure, testable, composable
export function calculateDiscount(
  price: number,
  discountPercent: number,
): number {
  if (price < 0 || discountPercent < 0 || discountPercent > 100) {
    throw new InvalidArgumentError("Invalid price or discount");
  }
  return price * (1 - discountPercent / 100);
}
```

**Standards**:

- No side effects
- Validate inputs
- Single return type
- Comprehensive tests

#### Service Methods

```typescript
// Good: Clear dependencies, error handling
async function createUserAccount(
  userData: CreateUserDto,
  dependencies: {
    userRepo: UserRepository;
    emailService: EmailService;
    logger: Logger;
  },
): Promise<Result<User, CreateUserError>> {
  const { userRepo, emailService, logger } = dependencies;

  try {
    // Validation
    const validation = validateUserData(userData);
    if (validation.isError()) {
      return Err(validation.error);
    }

    // Business logic
    const user = await userRepo.create(userData);
    await emailService.sendWelcome(user.email);

    logger.info("User created", { userId: user.id });
    return Ok(user);
  } catch (error) {
    logger.error("Failed to create user", error);
    return Err(new CreateUserError("Creation failed", error));
  }
}
```

#### Event Handlers

```typescript
// Good: Focused, async, error boundary
async function handleUserLogin(
  event: LoginEvent,
  context: HandlerContext,
): Promise<void> {
  const { userId, timestamp } = event;

  try {
    await validateSession(userId);
    await updateLastLogin(userId, timestamp);
    await notifySecurityTeam(userId);

    context.logger.info("Login handled", { userId });
  } catch (error) {
    context.logger.error("Login handler failed", { userId, error });
    await context.deadLetter.send(event);
  }
}
```

### 4. Refactoring Recommendations

**Extract Method Pattern**:

```typescript
// Before: Complex function
function processOrder(order: Order): ProcessedOrder {
  // 100 lines of mixed validation, calculation, and persistence
}

// After: Composed functions
function processOrder(order: Order): ProcessedOrder {
  const validated = validateOrder(order);
  const priced = calculatePricing(validated);
  const taxed = applyTaxes(priced);
  return persistOrder(taxed);
}
```

**Parameter Object Pattern**:

```typescript
// Before: Too many parameters
function createInvoice(
  customerId: string,
  items: Item[],
  discount: number,
  taxRate: number,
  dueDate: Date,
  notes: string,
): Invoice;

// After: Parameter object
interface CreateInvoiceParams {
  customerId: string;
  items: Item[];
  discount?: number;
  taxRate: number;
  dueDate: Date;
  notes?: string;
}

function createInvoice(params: CreateInvoiceParams): Invoice;
```

### 5. Integration with QDIRECTOR

**Function Validation Results**:

```yaml
function_validation:
  function_name: "processPayment"
  file_path: "/src/services/PaymentService.ts:45"
  validation_score: 85%
  complexity_score: 12

  issues:
    must_fix:
      - type: "missing_error_handling"
        details: "No catch block for payment gateway errors"
        line: 67

    should_fix:
      - type: "high_complexity"
        details: "Cyclomatic complexity: 12 (target: 10)"
        suggestion: "Extract payment validation logic"

    consider:
      - type: "parameter_count"
        details: "4 parameters, consider parameter object"

  recommendation: "refactor"

  refactor_plan:
    - "Extract validatePaymentData() method"
    - "Create PaymentGatewayError class"
    - "Add retry logic with exponential backoff"
```

### 6. Function Quality Metrics

**Track Over Time**:

```yaml
function_metrics:
  average_complexity: 6.3
  average_line_count: 28
  pure_function_ratio: 0.42
  test_coverage: 0.89

  improvement_trends:
    complexity: -15% # Reduced over sprint
    coverage: +8% # Increased coverage

  top_complex_functions:
    - "OrderProcessor.processOrder": 18
    - "UserService.syncPermissions": 15
    - "ReportGenerator.buildReport": 14
```

## Usage Examples

```bash
# Validate specific function
/qcheckf processPayment

# Validate all functions in file
/qcheckf --file /src/services/PaymentService.ts

# Validate with specific standards
/qcheckf --strict --max-complexity 8
```

This focused validation ensures each function is a well-crafted, maintainable
unit of code.
