// WARNING: Do not manually modify this file. It was generated using:
// https://github.com/dillonkearns/elm-typescript-interop
// Type definitions for Elm ports
export as namespace Elm


export interface App {
  ports: {
    hello: {
      subscribe(callback: (data: string) => void): void
    }
    reply: {
      send(data: number): void
    }
  }
}
    

export namespace Main {
  export function fullscreen(): App
  export function embed(node: HTMLElement | null): App
}