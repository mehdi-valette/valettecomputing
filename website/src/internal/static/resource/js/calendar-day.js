import { minutesToHours } from "./minutes-to-hours.js";
import { VsPeriod } from "./period.js";

export class VsCalendarDay extends HTMLElement {
  /** @type {HTMLElement} main container of this custom element */
  #container;

  /** @type {HTMLElement} main container of this custom element */
  #graphContainer;

  /** @type {number} */
  #totalMinutes;

  constructor() {
    super();

    this.#container = this.#createContainer();
    this.#graphContainer = this.#createGraphContainer();
    this.#totalMinutes = 1440;
  }

  connectedCallback() {
    this.attachShadow({ mode: "open" }).append(this.#container);

    this.#container.append(this.#createTimeline());
    this.#container.append(this.#graphContainer);
    this.#container.append(document.createElement("div"));
  }

  // get the number of pixels per minute
  get pixelStep() {
    return this.#container.scrollHeight / this.#totalMinutes;
  }

  get totalMinutes() {
    return this.#totalMinutes;
  }

  get height() {
    return this.#container.clientHeight;
  }

  get top() {
    return (
      this.#container.getBoundingClientRect().top +
      document.documentElement.scrollTop
    );
  }

  #createContainer = () => {
    const container = document.createElement("div");
    container.style.width = "500px";
    container.style.height = "1200px";
    container.style.display = "grid";
    container.style.gridTemplateColumns = "1fr 1fr 1fr";

    return container;
  };

  #createGraphContainer() {
    const graphContainer = document.createElement("div");
    graphContainer.style.width = "500px";
    graphContainer.style.height = "100%";
    graphContainer.style.position = "relative";
    graphContainer.style.backgroundColor = "lightgray";

    const slot = document.createElement("slot");
    slot.name = "graph";
    slot.addEventListener("slotchange", (evt) => {
      slot.assignedElements().forEach((element) => {
        if (element instanceof VsPeriod) element.setParent(this);
      });
    });
    graphContainer.append(slot);

    return graphContainer;
  }

  #createTimeline = () => {
    const timelineContainer = document.createElement("div");
    timelineContainer.style.width = "10rem";
    timelineContainer.style.position = "relative";

    for (let i = 0; i < 48; i++) {
      const line = document.createElement("div");
      line.style.top = this.pixelStep * 30 * i + "px";
      line.style.width = "100%";
      line.append(document.createTextNode(minutesToHours(i * 30)));
      line.style.borderTopWidth = "1px";
      line.style.borderTopColor = "blue";
      line.style.borderTopStyle = "solid";
      line.style.position = "absolute";

      timelineContainer.append(line);
    }

    return timelineContainer;
  };
}

customElements.define("vs-calendar-day", VsCalendarDay);
