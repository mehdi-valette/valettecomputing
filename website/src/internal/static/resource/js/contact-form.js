function handleContactForm() {
  const contactForm = document.getElementById("contact-form");

  if (!(contactForm instanceof HTMLFormElement)) return false;

  const form = new FormData(contactForm);

  contactForm.getElementsByTagName("button")[0].disabled = true;

  sendContactForm(form);

  return false;
}

/** @param {FormData} form  */
async function sendContactForm(form) {
  const contactForm = document.getElementById("contact-form");
  if (!(contactForm instanceof HTMLFormElement)) {
    console.error("the contact form is not found");
    return;
  }

  const container = document.getElementById("contact-form-container");
  if (container == null) {
    console.error("the contact form container is not found");
    return;
  }

  try {
    const response = await fetch(contactForm.action, {
      body: new URLSearchParams(form).toString(),
      method: "post",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
    });

    const data = await response.text();
    const node = new DOMParser()
      .parseFromString(data, "text/html")
      .body.children.item(0);
    container.appendChild(node);

    animationObserver.observe(node);

    contactForm.getElementsByTagName("button")[0].disabled = false;
  } catch (e) {
    console.log(e);
    container.innerHTML = "oups";
  }
}

async function removeContactSuccess() {
  const element = document.getElementById("contact-form-success");

  if (element == null) return;

  element.classList.remove("intersecting");

  element.ontransitionend = () => {
    animationObserver.unobserve(element);
    element.remove();
  };
}
