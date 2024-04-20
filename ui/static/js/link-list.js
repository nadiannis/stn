const linkDeleteButtons = document.querySelectorAll('.btn-link-delete');

linkDeleteButtons.forEach((linkDeleteButton) => {
  const id = linkDeleteButton.dataset.id;
  linkDeleteButton.addEventListener('click', () => deleteLink(id));
});
