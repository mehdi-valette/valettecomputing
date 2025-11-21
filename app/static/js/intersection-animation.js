document.addEventListener("DOMContentLoaded", () => {
  const observer = new IntersectionObserver(
    entries => {
      entries.forEach(entry => {
        console.log(entry);
        if (entry.isIntersecting) {
          entry.target.classList.add("intersecting");
        }
      })
    }
  );

  const test = document.querySelectorAll(".animatable");
  console.log(test);
  test.forEach(el => observer.observe(el));
}, { once: true })