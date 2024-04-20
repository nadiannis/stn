lucide.createIcons();

const backButton = document.querySelector('.btn-back');

if (backButton) {
  backButton.addEventListener('click', () => {
    history.back();
  });
}
