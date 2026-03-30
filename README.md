# TRAIN TRACK

**Train Track** is a self-hosted application for tracking weight lifting sessions. It's built with Go, HTML, vanilla JS, and uses SQLite for data storage.

![Screenshot](https://github.com/section14/train-track/blob/main/screenshot.jpg)

## Installation

**Requires Go v1.25.0+**

```git clone git@github.com:section14/train-track.git```

```go build -o train-track .```

Change `localhost` to your desired host, and `8080` to the port you want to run the application on

```mkdir .data && cp test-db/db.sqlite .data```

```./train-track prod localhost 8080```

---

**Note:** If you intend to run the application outside of the repo directory, you'll need to copy the `/static`, `/extracted`, `/templates` and `/.data` directories along with it. Embedding `/static` and `/templates` in the binary is on the roadmap.

---

## Development

Outside of a light, self-hosted app to track workouts with, I was motivated to build a dynamically rendered web app without any frontend frameworks or dependencies. **Train Track** is a multi-page app, which uses a mix of server rendered HTML partials (Go Templates) and Web Components to build out each page. I'll outline the techniques used in the following sections.

**Note:** This app uses [cross-page View Transitions](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API) and the [Popover API's hint](https://developer.mozilla.org/en-US/docs/Web/API/Popover_API/Using#using_hint_popover_state). These aren't fully implemented in Firefox and Safari, respectively. So, **Train Track** only renders 100% correctly in Chromium based browsers.

### Partials

The partials are constructed in a similar way to Web Components or React components. There's an HTML section and a Javascript section, which is incased in `<script>` tags. Here's a simple, stripped down example:

```html
<div class="exercise-list p-3">
    <ul>
        {{range .}}
        <li id="item-row-{{.ID}}">
            <div id="item-name-{{.ID}}" class="item-name">
                {{.Name}}
            </div>
            <button data-delete="{{.ID}}" type="button" class="btn-edit">
                <svg-icon icon="edit"></svg-icon>
            </button>
        </li>
        {{end}}
    </ul>
</div>

<script>
    import { Delete } from '/static/js/request.js'
    import { swapContent } from '/static/js/swap.js'

    //@handle:data-delete
    export const deleteExercise = (id) => {
        Delete(`/api/exercises/${id}`).then(() => {
            swapContent("/api/partials/exercises", "exercise-list-container")
        }).catch((err) => {
            console.error("delete error: ", err)
        })
    }
</script>
```

During the application's build stage, all HTML files have the Javascript extracted, and placed into a new file of the same name (eg. exercise-list.html -> exercise-list.js). Which is then imported on pages that will use that partial. This lets us develop the Javascript along side HTML, but only have a singe Js file for n number of HTML partials.

```
//@handle:data-delete
```

This is a directive which tells the parser to associate a particular function (deleteExercise in this case) with a click Event + data-* attribute. "click" event handlers are generated at build time, and automatically imported into a page.

In order to prevent naming collisions, and `@handle` referenced data-* attributes are modified with an id, unique to that HTML partial. So, `data-delete` will be changed to something like `data-delete-123456789` when the app is launched.

`swapContent` is a utility function which issues a `GET` request to the first argument, and replaces the content of the second argument (a div, typically) with the result.

---

### Auto-loading

Automatically loading partials is more difficult than it might seem. I looked at other libraries like HTMX for inspiration, but ultimately wrote a (Light DOM) Web Component to handle this. Which worked out well.

```html
<partial-loader
    target-id="workout-1"
    endpoint="/api/partials/workouts/1"
>
```

This is basically a wrapper around `swapContent()` which injects the result of `endpoint` into `target-id` when the component is mounted.

---

### Web Components - Pros

Comparing Web Components and React has been interesting. While the approaches are different, I've found Web Components to be equal in functionality on most counts, and preferred in some instances. Using `connectedCallback` and `disconnectedCallback` in place of `useEffect` is more direct (although React < 18 basically had this).

Watching attributes (ie. props) is more or less the same. Plus, Web Components give you the old value on change, which is a nice bonus. Another highlight is the ability to update component state synchronously. React's async-only approach is often something that makes complex state updates a problem. Even with `useReducer`.

--- 

### Web Components - Cons

The primary pain points for Web Components are including them on a page and styling. If you have an HTML file for the `<template>` and a Javascript file for the component class as I do, you have to supply both files to a page. While you can write the HTML in a template literal, it fells messy to me. Go templating makes file inclusion relatively trivial, but it's more work than a React import.

Styling leaves quite a bit to be desired. While CSS variables, and certain properties like font-family will cross into the Shadow DOM, none of your classes will. So, utility classes can't be used in them. There are ways to define styles from outside of a Web Component by using the [part](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Selectors/::part) pseudo selector, but your left defining styles twice. Which isn't ideal.

In my opinion, an option to annotate existing CSS classes, exposing them to the Shadow DOM, would work best.

---

## Conclusion

While I'm not convinced the system I've developed for this app is 100% viable, it's a worthy exploration into alternatives to the status quo in web development. And, I believe it exposes some techniques for building a dynamic UI without a frontend framework. One thing that's for sure: Javascript and CSS have come a long way by providing tools that were previously only available via 3rd party frameworks and libraries. There's an opportunity to massively reduce dependencies, and decouple apps from the endless churn of frontend libraries. Those goals alone make projects like this worth the effort.

Ultimately, the trade offs between React components and server rendered HTML partials are different. For example, I disable adding new exercises to a workout, and deleting a workout when modifying an exercise's properties (eg. sets, reps). As adding a new exercise will pull unsaved data from API, discarding changes. This is less of an issue with something like React.

## Todo

- [ ] Store the exercise list (in a `<select>`) a level up to prevent redundant API calls
- [ ] Provide a confirmation Popover when deleting a workout and exercise (from the main list)
- [ ] Reporting to track progress on exercises
- [ ] Embed `/templates` and `/static` into the Go binary for production builds
