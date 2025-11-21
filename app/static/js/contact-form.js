function handleContactForm() {
  const contactForm = document.getElementById("contact-form")
  const form = new FormData(contactForm)

  contactForm.getElementsByTagName("button")[0].disabled = true

  sendContactForm(form)

  return false;
}

async function sendContactForm(form) {
  const contactForm = document.getElementById("contact-form")
  const container = document.getElementById("contact-form-container")

  try {
    const response = await fetch("/contact", {
      body: new URLSearchParams(form).toString(),
      method: "post",
      headers: { "Content-Type": "application/x-www-form-urlencoded" }
    })

    const data = await response.text()
    const node = new DOMParser().parseFromString(data, "text/html").body.children.item(0)
    container.appendChild(node)

    animationObserver.observe(node)

    contactForm.getElementsByTagName("button")[0].disabled = false
  } catch (e) {
    console.log(e)
    container.innerHTML = "oups"
  }
}

async function removeContactSuccess() {
  const element = document.getElementById("contact-form-success")

  if (element == null) return

  element.classList.remove("intersecting")

  element.ontransitionend = () => {
    animationObserver.unobserve(element)
    element.remove()
  }
}