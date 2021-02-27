export function moveElement<T>(oldList: Array<T>, sourceIdx: number, destIdx: number): Array<T> {
  const newList = [...oldList];
  const movedElement = newList.splice(sourceIdx, 1)[0];
  newList.splice(destIdx, 0, movedElement);
  return newList;
}

export function removeOneFromList<T>(list: Array<T>, idx: number): Array<T> {
  return [
    ...list.slice(0, idx),
    ...list.slice(idx+1, list.length)
  ];
}

export function addOneToList<T>(list: Array<T>, idx: number, element: T) {
  return [
    ...list.slice(0, idx),
    element,
    ...list.slice(idx, list.length)
  ];
}
