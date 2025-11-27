class VsCalendarDay extends HTMLElement {
  constructor() {
    super();
    this.container = this.createContainer();
  }

  connectedCallback() {
    this.append(this.container);

    // Set up observer
    this.observer = new MutationObserver(this.onMutation);

    // Watch the Light DOM for child node changes
    this.observer.observe(this, {
      childList: true,
    });
  }

  disconnectedCallback() {
    this.observer?.disconnect();
  }

  /** @type {MutationCallback} */
  onMutation = (mutations) => {
    mutations.forEach((mutation) => {
      this.container.append(...mutation.addedNodes);

      mutation.addedNodes.forEach((node) => {
        if (!(node instanceof VsPeriod)) return;

        node.setParent(this.container);
      });
    });
  };

  createContainer() {
    const container = document.createElement("div");
    container.style.width = "500px";
    container.style.height = "1200px";
    container.style.position = "relative";
    container.style.backgroundColor = "lightgray";

    return container;
  }
}

customElements.define("vs-calendar-day", VsCalendarDay);
