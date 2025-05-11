document.getElementById("loginForm").addEventListener("submit", async function (e)  {
    e.preventDefault();

    const form = e.target;
    const username=form.username.value;
    const password = form.password.value;
    const nonce = form.nonce.value;

    // Step 1: SHA512(password)
    const passHash = await sha512(password);

    // Step 2: SHA512(passHash + nonce)
    const finalHash = await sha512(passHash + nonce);

    form.password.value = finalHash;
    const response = await fetch("/login_check", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username:username, password: finalHash, nonce:nonce })
      });
      if (response.status===200 ){
        const result = await response.json();
        if (result.token!==undefined && result.token) {
          localStorage.setItem("jwt", result.token);
          document.cookie = `token=${result.token}; path=/; SameSite=Strict`;
          window.location.href = "/home";
        }
        else
        {
          document.getElementById("responseMessage").innerHTML="Invalid UserName/Password";
          document.getElementById("password").value="";
        }
      } else {
        document.getElementById("responseMessage").innerHTML="Invalid UserName/Password";
        document.getElementById("password").value="";
      }
});

async function sha512(str) {
    const buffer = new TextEncoder().encode(str);
    const digest = await crypto.subtle.digest("SHA-512", buffer);
    return Array.from(new Uint8Array(digest)).map(b => b.toString(16).padStart(2, "0")).join("");
}