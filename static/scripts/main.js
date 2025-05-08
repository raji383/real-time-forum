document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const signupForm = document.getElementById('loginForm')

  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    signupForm.style.display='none'

  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    signupForm.style.display='block'
  });
});
