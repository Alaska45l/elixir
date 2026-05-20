export function reveal(node: Element) {
  const observer = new IntersectionObserver(
    (entries) => entries.forEach((entry) => entry.isIntersecting && entry.target.classList.add('visible')),
    { threshold: 0.12 }
  );
  node.classList.add('reveal');
  observer.observe(node);
  return { destroy: () => observer.disconnect() };
}
