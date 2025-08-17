# Execution Flow - How Data Moves Through the System

## 🔄 **Overview**

This document explains how data flows through TokenWatch CLI for different types of commands. Understanding these flows is crucial for debugging, optimization, and adding new features.

## 📊 **Flow Diagrams**

### **1. Single Platform Command Flow (e.g., `tokenwatch openai`)**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │ openai_cmd  │    │ openai.go   │    │ OpenAI API  │
│             │    │             │    │             │    │             │
│ tokenwatch  │───▶│ 1. Parse    │───▶│ 1. Get API  │───▶│ 1. Usage    │
│   openai    │    │   flags     │    │   key       │    │   endpoint  │
│             │    │ 2. Validate │    │ 2. Build    │    │ 2. Cost     │
│             │    │   period    │    │   request   │    │   endpoint  │
│             │    │ 3. Call     │    │ 3. Make     │    │ 3. Return   │
│             │    │   provider  │    │   HTTP call │    │   JSON data │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                           │                   │                   │
                           ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │ openai_cmd  │    │ openai.go   │    │             │
│             │    │             │    │             │    │             │
│             │◀───│ 1. Receive  │◀───│ 1. Parse    │◀───│             │
│             │    │   data      │    │   response  │    │             │
│             │    │ 2. Format   │    │ 2. Transform│    │             │
│             │    │   output    │    │   to models │    │             │
│             │    │ 3. Display  │    │ 3. Return   │    │             │
│             │    │   table     │    │   models    │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### **2. Multi-Platform Command Flow (e.g., `tokenwatch all`)**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │  all_cmd    │    │  Platform   │
│             │    │             │    │  Detection  │
│ tokenwatch  │───▶│ 1. Get      │───▶│ 1. Check    │
│    all      │    │   platforms │    │   config    │
│             │    │ 2. Launch   │    │ 2. Return   │
│             │    │   parallel  │    │   list      │
│             │    │   workers   │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                PARALLEL EXECUTION                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ Goroutine 1 │  │ Goroutine 2 │  │ Goroutine 3 │     │
│  │             │  │             │  │             │     │
│  │ openai.go   │  │ grok.go     │  │ anthropic.go│     │
│  │ GetConsumptionSummary()      │                 │     │
│  │ GetPricingSummary()          │                 │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
│           │              │              │               │
│           ▼              ▼              ▼               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ OpenAI API  │  │ Grok API    │  │ Anthropic   │     │
│  │             │  │             │  │ API         │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │  all_cmd    │    │  Results    │
│             │    │             │    │  Collection │
│             │◀───│ 1. Collect  │◀───│ 1. Channel  │
│             │    │   results   │    │   results   │
│             │    │ 2. Aggregate│    │ 2. Error    │
│             │    │   data      │    │   handling  │
│             │    │ 3. Display  │    │ 3. Success  │
│             │    │   combined  │    │   data      │
└─────────────┘    └─────────────┘    └─────────────┘
```

### **3. Setup Command Flow (e.g., `tokenwatch setup`)**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │ setup_cmd   │    │  Platform   │
│             │    │             │    │  Selection  │
│ tokenwatch  │───▶│ 1. Show     │───▶│ 1. Display │
│   setup     │    │   platforms │    │   options   │
│             │    │ 2. Get      │    │ 2. Get      │
│             │    │   choice    │    │   user      │
│             │    │ 3. Validate │    │   input     │
│             │    │   selection │    │ 3. Validate │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │ setup_cmd   │    │  Config     │
│             │    │             │    │  Manager    │
│             │◀───│ 1. Prompt   │◀───│ 1. Load     │
│             │    │   for API   │    │   existing  │
│             │    │   key       │    │   config    │
│             │    │ 2. Masked   │    │ 2. Create   │
│             │    │   input     │    │   new if    │
│             │    │ 3. Save     │    │   needed    │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    User     │    │ setup_cmd   │    │  File       │
│             │    │             │    │  System     │
│             │◀───│ 1. Confirm  │◀───│ 1. Write    │
│             │    │   success   │    │   config    │
│             │    │ 2. Show     │    │ 2. Set      │
│             │    │   next      │    │   permissions│
│             │    │   steps     │    │ 3. Verify   │
└─────────────┘    └─────────────┘    └─────────────┘
```

## 🔍 **Detailed Flow Analysis**

### **Single Platform Command Flow Steps**

#### **Phase 1: Command Parsing & Validation**
1. **User Input**: `tokenwatch openai --period 30d`
2. **Flag Parsing**: Extract `--period` flag value
3. **Validation**: Ensure period is valid (7d, 30d, 90d)
4. **Default Handling**: Set default period if none specified

#### **Phase 2: Configuration Loading**
1. **Config Init**: Load configuration from `~/.tokenwatch/config.yaml`
2. **API Key Retrieval**: Get OpenAI API key from config
3. **Validation**: Ensure API key exists and is valid

#### **Phase 3: Provider Creation**
1. **Factory Call**: Call `getProvider("openai")`
2. **Provider Creation**: Create new OpenAI provider instance
3. **Availability Check**: Verify provider is properly configured

#### **Phase 4: Data Retrieval**
1. **API Call**: Make HTTP request to OpenAI usage endpoint
2. **Response Handling**: Parse JSON response from OpenAI
3. **Error Handling**: Handle API errors, rate limits, timeouts
4. **Data Transformation**: Convert OpenAI format to common models

#### **Phase 5: Data Processing**
1. **Aggregation**: Calculate totals across all models
2. **Period Filtering**: Filter data to requested time period
3. **Calculations**: Compute averages, costs per token, etc.

#### **Phase 6: Output Generation**
1. **Table Creation**: Build formatted table with data
2. **Styling**: Apply colors, borders, alignment
3. **Display**: Output to terminal

### **Multi-Platform Command Flow Steps**

#### **Phase 1: Platform Detection**
1. **Config Scan**: Check all supported platforms for API keys
2. **Platform List**: Create list of configured platforms
3. **Provider Creation**: Create provider instances for each platform

#### **Phase 2: Parallel Execution Setup**
1. **Worker Creation**: Launch goroutine for each platform
2. **Channel Setup**: Create result channels for data collection
3. **WaitGroup Setup**: Track completion of all workers

#### **Phase 3: Concurrent API Calls**
1. **Parallel Execution**: Each platform makes API calls simultaneously
2. **Independent Processing**: Each provider handles its own API logic
3. **Result Collection**: Results sent through channels

#### **Phase 4: Result Aggregation**
1. **Data Collection**: Collect results from all channels
2. **Error Handling**: Handle failed platform calls gracefully
3. **Data Combination**: Aggregate data across all platforms

#### **Phase 5: Combined Display**
1. **Unified Dashboard**: Show combined metrics across platforms
2. **Platform Breakdown**: Display individual platform contributions
3. **Performance Metrics**: Show parallel execution benefits

## ⚡ **Performance Characteristics**

### **Single Platform Commands**
- **Response Time**: 100-500ms (depending on API response time)
- **Memory Usage**: Low (single provider instance)
- **Network Calls**: 2-4 API calls per command
- **Scalability**: Linear with API response time

### **Multi-Platform Commands**
- **Response Time**: Max(platform_response_times) + 50ms overhead
- **Memory Usage**: Medium (multiple provider instances)
- **Network Calls**: 2-4 API calls per platform (parallel)
- **Scalability**: Excellent (parallel execution)

### **Performance Optimization**
- **Parallel Execution**: Reduces total response time
- **Connection Pooling**: Reuse HTTP connections
- **Caching**: Cache frequently accessed data
- **Timeout Management**: Prevent hanging on slow APIs

## 🚨 **Error Handling Flows**

### **API Failure Handling**
```
API Call Fails → Provider Returns Error → Command Handles Error → User Sees Message
```

### **Configuration Error Handling**
```
Config Load Fails → Fallback to Defaults → User Sees Warning → Setup Command Suggested
```

### **Platform Unavailability**
```
Platform Down → Provider Returns Error → Command Continues → Other Platforms Still Work
```

## 🔄 **Data Transformation Flows**

### **OpenAI API → Common Models**
```
OpenAI JSON → openai.go → models.ConsumptionSummary → openai_cmd.go → Table Display
```

### **Grok API → Common Models**
```
Grok JSON → grok.go → models.ConsumptionSummary → grok_cmd.go → Table Display
```

### **Multi-Platform Aggregation**
```
[OpenAI Models, Grok Models, ...] → all_cmd.go → Combined Display
```

## 📈 **Monitoring & Debugging Flows**

### **Debug Mode Flow**
```
Debug Flag Set → Verbose Logging → API Call Details → Response Data → User Debug Info
```

### **Performance Monitoring**
```
Command Start → Timer Start → API Calls → Timer End → Performance Metrics Display
```

---

**Understanding these execution flows is essential for maintaining, debugging, and extending TokenWatch CLI. Each flow follows the same architectural principles while optimizing for the specific use case.**
