document.addEventListener("DOMContentLoaded", function () {
    // Login Form Handling
    const loginForm = document.getElementById("loginForm");
    if (loginForm) {
        loginForm.addEventListener("submit", async function (e) {
            e.preventDefault();

            const username = loginForm.username.value;
            const password = loginForm.password.value;
            const nonce = loginForm.nonce.value;
            const role = loginForm.role.value;

            const passHash = await sha512(password);
            const finalHash = await sha512(passHash + nonce);
            loginForm.password.value = finalHash;

            const response = await fetch("/login_check", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password: finalHash, nonce, role })
            });

            if (response.status === 200) {
                const result = await response.json();
                if (result.token) {
                    localStorage.setItem("jwt", result.token);
                    document.cookie = `token=${result.token}; path=/; SameSite=Strict`;

                    // ✅ Use server-defined redirect path
                    const redirectPath = result.redirect_path || "/student_dashboard";
                    window.location.href = redirectPath;
                } else {
                    showError();
                }
            } else {
                showError();
            }

            function showError() {
                const responseMsg = document.getElementById("responseMessage");
                if (responseMsg) {
                    responseMsg.innerHTML = "Invalid Username/Password";
                }
                loginForm.password.value = "";
            }
        });
    }

    // Register Form Handling
    const registerForm = document.getElementById("registerForm");
    if (registerForm) {
        registerForm.addEventListener("submit", async function (e) {
            e.preventDefault();

            const passwordInput = registerForm.querySelector("#password");
            const password = passwordInput.value;

            const hashedPassword = await sha512(password);
            passwordInput.value = hashedPassword;

            registerForm.submit(); // manually submit after hashing
        });
    }
});

// SHA-512 hashing function
async function sha512(str) {
    const buffer = new TextEncoder().encode(str);
    const digest = await crypto.subtle.digest("SHA-512", buffer);
    return Array.from(new Uint8Array(digest))
        .map(b => b.toString(16).padStart(2, "0"))
        .join("");
}
