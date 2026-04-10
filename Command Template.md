# 🚀 Custom Command Templates: Extend XSH to Any Remote Tool

> **Transform XSH into a universal remote connection manager.** Support SSH, Azure CLI, AWS SSM, kubectl exec, and any CLI tool with a single, elegant abstraction.

## 🎯 The Universal Connection Problem

Every remote connection tool shares the same core elements:

| Component | Description |
|-----------|-------------|
| **Address** | DNS name or IP address of the target machine |
| **Port** | Remote port for connection initiation |
| **User** | Authentication username |
| **Identity** | Path to private key or authentication method |
| **Flags** | Tool-specific command-line arguments |

XSH already manages these properties beautifully for SSH. **But what if we told you XSH can manage connections for ANY remote tool?**

## ✨ Introducing Custom Command Templates

Custom Command Templates transform XSH from an SSH manager into a **universal remote connection orchestrator**. With a simple template system, XSH can now generate connection strings for any CLI tool in the world.

### How It Works

XSH introduces a `tools` table with two powerful columns:

- **`name`**: Identifier for your CLI tool (e.g., "azure-cli", "kubectl", "custom-ssh")
- **`template`**: Connection string template with variable placeholders

### 🔧 Template Variables

Use these variables in your templates:

```bash
${address}           # Target hostname or IP
${port}              # Connection port
${user}              # Username
${identity_file_path} # Path to private key
${extra_flags}       # Additional CLI flags
```

## 📖 Real-World Examples

### Example 1: Custom SSH Tool

**Template Definition:**
```bash
myssh -i ${identity_file_path} -p ${port} ${extra_flags} ${user}@${address}
```

**Host Configuration:**
```json
{
    "name": "production-server",
    "address": "prod.example.com",
    "port": 2026,
    "user": "deploy",
    "extra_flags": "-4 -A"
}
```

**Generated Command:**
```bash
xsh connect production-server
# Executes: myssh -i /path/to/key -p 2026 -4 -A deploy@prod.example.com
```

### Example 2: Azure CLI Integration

**Template Definition:**
```bash
az ssh vm -n ${address} ${extra_flags}
```

**Host Configuration:**
```json
{
    "name": "azure-vm-web",
    "address": "web-server-01",
    "extra_flags": "--resource-group production-rg --subscription prod-sub"
}
```

**Generated Command:**
```bash
xsh connect azure-vm-web
# Executes: az ssh vm -n web-server-01 --resource-group production-rg --subscription prod-sub
```

### Example 3: Kubernetes Pod Access

**Template Definition:**
```bash
kubectl exec -it ${address} ${extra_flags} -- /bin/bash
```

**Host Configuration:**
```json
{
    "name": "api-pod",
    "address": "api-deployment-7d4b8c9f-x4k2z",
    "extra_flags": "--namespace production --context prod-cluster"
}
```

**Generated Command:**
```bash
xsh connect api-pod
# Executes: kubectl exec -it api-deployment-7d4b8c9f-x4k2z --namespace production --context prod-cluster -- /bin/bash
```

## 🎨 Template Flexibility

### Selective Variable Usage

Not every template needs every variable. Create focused templates for specific use cases:

```bash
# Minimal template
docker exec -it ${address} /bin/bash

# Port-specific template  
telnet ${address} ${port}

# Identity-only template
ssh-copy-id -i ${identity_file_path} ${user}@${address}
```

### Advanced Template Patterns

```bash
# Jump host with custom SSH
ssh -i ${identity_file_path} -J ${extra_flags} ${user}@${address}

# AWS SSM Session Manager
aws ssm start-session --target ${address} ${extra_flags}

# Custom debugging tool
mytool --host ${address} --port ${port} --debug ${extra_flags}
```

## 🌟 Benefits of Custom Templates

### 🔄 **Universal Abstraction**
- One interface for all remote tools
- Consistent command patterns across your infrastructure
- Reduces cognitive load when switching between tools

### 📚 **Knowledge Preservation**
- Store complex command patterns once, use everywhere
- Share tool configurations across teams
- Document connection requirements in a structured format

### ⚡ **Instant Productivity**
- No more memorizing tool-specific syntax
- Quick switching between different connection methods
- Batch operations across multiple tools

### 🔧 **Infinite Extensibility**
- Support any CLI tool with connection capabilities
- Create organization-specific connection standards
- Adapt to new tools without changing workflows

## 🚀 Getting Started

1. **Define your template:**
   ```bash
   xsh put tool
   ```

2. **Configure your hosts** with the tool-specific properties

3. **Connect seamlessly:**
   ```bash
   xsh connect my-host
   ```

---

**Custom Command Templates unlock XSH's true potential** — transforming it from an SSH manager into your organization's universal remote connection hub. Your imagination is the only limit.