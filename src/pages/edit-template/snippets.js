// @flow

export const SNIPPET_OPTIONS = [
  {
    id: "first_name",
    label: "First Name",
    snippet: "{{ .Contact.FirstName }}",
  },
  {
    id: "last_name",
    label: "Last Name",
    snippet: "{{ .Contact.LastName }}",
  },
  {
    id: "track_link",
    label: "Tracked Link",
    snippet:
      '<a target="_blank" href="{{trackURL "https://google.com"}}">My Link</a>',
  },
  {
    id: "unsubscribe_link",
    label: "Unsubscribe Link",
    snippet: '<a target="_blank" href="{{ .UnsubscribeURL }}">Unsubscribe</a>',
  },
  {
    id: "open_track_image",
    label: "Open Track Image",
    snippet: '<img src="{{ .OpenEventImageURL }}">',
  },
];
