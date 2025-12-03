/**
 * @param {number} mn
 * @returns {string}
 */
export function minutesToHours(mn) {
  const hours = Math.floor(mn / 60);
  const minutes = mn % 60;

  return (
    hours.toString().padStart(2, "0") +
    ":" +
    minutes.toString().padStart(2, "0")
  );
}
