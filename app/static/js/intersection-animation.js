const animationObserver = new IntersectionObserver(
  entries => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add("intersecting")
      }
    })
  }
);

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".animatable").forEach(el => animationObserver.observe(el))
}, { once: true })