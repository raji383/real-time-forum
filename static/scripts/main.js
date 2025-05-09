document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const loginForm = document.getElementById('loginForm');
  const mainPage = document.getElementById('mainPage');

  // Display signup form
  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    loginForm.style.display = 'none'; // Hide login form when sign up
  });

  // Display login form
  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    loginForm.style.display = 'block'; // Show login form when sign in
  });

  async function handleLoginSubmit(e) {
    e.preventDefault();
  
    const formData = new FormData(loginForm);
    const urlEncodedData = new URLSearchParams();
  
    // تحويل FormData إلى URLSearchParams
    for (let pair of formData.entries()) {
      urlEncodedData.append(pair[0], pair[1]);
    }
  
    try {
      const response = await fetch('/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',  
        },
        body: urlEncodedData.toString(),  
      });
  
      if (response.ok) {
        const result = await response.json();
  
        if (result.success) {
          container.style.display = 'none';
          mainPage.style.display = 'block';
  
          mainPage.innerHTML = `
            <h1>Welcome!</h1>
            <p>${result.message}</p>
            ${result.nickname ? `<p>Nickname: ${result.nickname}</p>` : ""}
          `;
        } else {
          alert("Login failed: " + result.message);
        }
      } else {
        const errorText = await response.text();
        alert("Error: " + errorText);
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("Something went wrong. Please try again.");
    }
  }
  
  loginForm.addEventListener('submit', handleLoginSubmit);
  
  
});
