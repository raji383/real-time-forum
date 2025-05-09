document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const signupForm = document.getElementById('loginForm')
  const sign = document.getElementById('sing');
  const mainPage = document.getElementById('mainPage');


  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    signupForm.style.display='none'

  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    signupForm.style.display='block'
  });
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault(); 

    const formData = new FormData(loginForm);

    const response = await fetch('/', {
      method: 'POST',
      body: formData,
    });

    if (response.ok) {
      const result = await response.json();

      if (result.success) {
        console.log(result);
        container.style.display = 'none';
        mainPage.style.display = 'block';
      } else {
        alert("Login failed: " + result.message);
      }
    } else {
      const errorText = await response.text();
      alert("Error: " + errorText);
    }
  });
    
});
