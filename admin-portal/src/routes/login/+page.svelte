<script>
    let email = "";
    let password = "";
    let isLoading = false;
    let errorMessage = "";

    async function handleLogin() {
        isLoading = true;
        errorMessage = "";

        try {
            await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate checking

            if (email === "test@123" && password === "1234") {
                // Navigate to dashboard on successful login
                window.location.href = "/";
            } else {
                errorMessage =
                    "Invalid credentials. Please check your email and password.";
            }
        } catch (e) {
            errorMessage = "An unexpected error occurred. Please try again.";
        } finally {
            isLoading = false;
        }
    }
</script>

<svelte:head>
    <title>Admin Login - Furniture Store</title>
</svelte:head>

<div class="login-wrapper">
    <div class="login-card">
        <div class="header">
            <span class="badge">🔐 Admin Portal</span>
            <h1>Welcome Back</h1>
            <p>Secure access for inventory management</p>
        </div>

        {#if errorMessage}
            <div class="error-alert">
                <span>⚠️</span>
                <span>{errorMessage}</span>
            </div>
        {/if}

        <form on:submit|preventDefault={handleLogin}>
            <div class="form-group">
                <label for="email">Email Address</label>
                <input
                    type="email"
                    id="email"
                    bind:value={email}
                    placeholder="admin@example.com"
                    required
                />
            </div>

            <div class="form-group">
                <label for="password">Password</label>
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    placeholder="Enter your password"
                    required
                />
            </div>

            <button type="submit" class="submit-btn" disabled={isLoading}>
                {#if isLoading}
                    <span class="spinner"></span>
                    <span>Authenticating...</span>
                {:else}
                    <span>Sign In</span>
                    <span>→</span>
                {/if}
            </button>
        </form>

        <div class="footer">
            <a href="http://localhost:5000/login" class="back-link">
                ← Back to Client Portal
            </a>
        </div>
    </div>
</div>

<style>
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }

    .login-wrapper {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        padding: 20px;
    }

    .login-card {
        width: 100%;
        max-width: 420px;
        background: white;
        border-radius: 24px;
        padding: 3rem;
        box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
    }

    .header {
        text-align: center;
        margin-bottom: 2rem;
    }

    .badge {
        display: inline-block;
        padding: 8px 20px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        color: white;
        border-radius: 20px;
        font-size: 12px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        margin-bottom: 16px;
    }

    .header h1 {
        font-size: 32px;
        font-weight: 800;
        color: #1a1a2e;
        margin-bottom: 8px;
    }

    .header p {
        font-size: 15px;
        color: #666;
    }

    .error-alert {
        background: rgba(255, 82, 82, 0.1);
        border: 1px solid rgba(255, 82, 82, 0.3);
        border-radius: 12px;
        padding: 14px 16px;
        margin-bottom: 20px;
        display: flex;
        align-items: center;
        gap: 10px;
        color: #c53030;
        font-size: 14px;
    }

    .form-group {
        margin-bottom: 20px;
    }

    .form-group label {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: #1a1a2e;
        margin-bottom: 8px;
    }

    .form-group input {
        width: 100%;
        padding: 14px 16px;
        font-size: 15px;
        background: #f8f9fa;
        color: #1a1a2e;
        border: 2px solid #e5e7eb;
        border-radius: 12px;
        transition: all 0.3s ease;
    }

    .form-group input:focus {
        outline: none;
        border-color: #667eea;
        background: white;
        box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
    }

    .form-group input::placeholder {
        color: #9ca3af;
    }

    .submit-btn {
        width: 100%;
        padding: 16px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        color: white;
        border: none;
        border-radius: 12px;
        font-size: 16px;
        font-weight: 700;
        cursor: pointer;
        transition: all 0.3s ease;
        box-shadow: 0 8px 20px rgba(102, 126, 234, 0.35);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
    }

    .submit-btn:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 12px 28px rgba(102, 126, 234, 0.45);
    }

    .submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .spinner {
        width: 18px;
        height: 18px;
        border: 3px solid rgba(255, 255, 255, 0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .footer {
        text-align: center;
        margin-top: 24px;
        padding-top: 20px;
        border-top: 1px solid #e5e7eb;
    }

    .back-link {
        color: #667eea;
        text-decoration: none;
        font-size: 14px;
        font-weight: 600;
        transition: opacity 0.2s;
    }

    .back-link:hover {
        opacity: 0.8;
        text-decoration: underline;
    }
</style>
