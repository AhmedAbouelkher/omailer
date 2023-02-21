package omailer

import (
	"fmt"
	"strings"
)

type Elem string

type HTML struct {
	elements []Elem
}

func NewHTML() *HTML {
	return &HTML{}
}

func (b *HTML) String() string {
	var sb strings.Builder
	for _, s := range b.elements {
		sb.WriteString(string(s))
	}
	return sb.String()
}

func (b *HTML) AddElem(elem ...Elem) *HTML {
	b.elements = append(b.elements, elem...)
	return b
}

func P(elem ...Elem) Elem {
	ss := mergeElems(elem...)
	s := fmt.Sprintf(
		`<p style="line-height: %s; font-size: 14px;">
		<span style="font-family: 'Open Sans', sans-serif; font-size: 14px; line-height: 19.6px;">
		%s
		</span></p>`,
		"140%",
		ss,
	)
	return Elem(s)
}

func Strong(elem ...Elem) Elem {
	ss := mergeElems(elem...)
	s := fmt.Sprintf(
		`<strong style="font-family: 'Open Sans', sans-serif; font-size: 14px; line-height: 19.6px;">
		%s
		</strong>`,
		ss,
	)
	return Elem(s)
}

func Padding(elem Elem, top, bottom int) Elem {
	s := fmt.Sprintf(
		`<div style="padding-top: %dpx; padding-bottom: %dpx;">%s</div>`,
		top,
		bottom,
		elem,
	)
	return Elem(s)
}

func A(elem Elem, link string) Elem {
	s := fmt.Sprintf(
		"<a href=\"%s\" target=\"_blank\">%s</a>",
		link,
		elem,
	)
	return Elem(s)
}

func List(elem ...Elem) Elem {
	var sb strings.Builder
	for _, e := range elem {
		ss := Span(e)
		sb.WriteString(fmt.Sprintf(`<li>%s</li>`, ss))
	}
	s := fmt.Sprintf(
		`<ul>
		%s
		</ul>`,
		sb.String(),
	)
	return Elem(s)
}

func Span(elem ...Elem) Elem {
	ss := mergeElems(elem...)
	s := fmt.Sprintf(
		`<span style="font-family: 'Open Sans', sans-serif; font-size: 14px; line-height: 19.6px;">
		%s
		</span>`,
		ss,
	)
	return Elem(s)
}

func Space(h int, w int) Elem {
	s := fmt.Sprintf(
		`<div style="height: %dpx; width: %dpx;"></div>`,
		h,
		w,
	)
	return Elem(s)
}

func Btn(elem Elem, link string) Elem {
	styles := strings.Join([]string{
		"text-decoration:none",
		"color:#ffffff",
		"text-align:center",
		"display:block",
		"border-radius:5px",
		"background-color:#0282a6",
		"padding-top:8px",
		"padding-right:16px",
		"padding-bottom:8px",
		"padding-left:16px",
	}, ";")
	s := fmt.Sprintf(
		`<a href="%s" style="%s" target="_blank">
		%s
		</a>`,
		link,
		styles,
		elem,
	)
	return Elem(s)
}

func Center(elem ...Elem) Elem {
	ss := mergeElems(elem...)
	s := fmt.Sprintf(
		`<table role="presentation" border="0" cellpadding="0" cellspacing="0" style="margin:0 auto">
		<tbody>
		  <tr>
			<td style="font-family: 'Open Sans', sans-serif;font-size:14px;">
			%s
			</td>
		  </tr>
	  </tbody>
	  </table>`,
		ss,
	)
	return Elem(s)
}

type TextStyle struct {
	Color      string
	FontSize   float32
	LineHeight float32
	FontWeight string
	Decoration string
}

func Text(elem Elem, s *TextStyle) Elem {
	df := &TextStyle{
		Color:      "#000000",
		FontSize:   14,
		LineHeight: 19.6,
		FontWeight: "normal",
		Decoration: "none",
	}
	var sb strings.Builder
	sb.WriteString(`<span style="font-family: 'Open Sans', sans-serif;`)
	if s.FontSize != 0 {
		sb.WriteString(fmt.Sprintf(` font-size: %.2fpx;`, s.FontSize))
	} else {
		sb.WriteString(fmt.Sprintf(` font-size: %.2fpx;`, df.FontSize))
	}
	if s.LineHeight != 0 {
		sb.WriteString(fmt.Sprintf(` line-height: %.2fpx;`, s.LineHeight))
	} else {
		sb.WriteString(fmt.Sprintf(` line-height: %.2fpx;`, df.LineHeight))
	}
	if s.Color != "" {
		sb.WriteString(fmt.Sprintf(` color: %s;`, s.Color))
	} else {
		sb.WriteString(fmt.Sprintf(` color: %s;`, df.Color))
	}
	if s.FontWeight != "" {
		sb.WriteString(fmt.Sprintf(` font-weight: %s;`, s.FontWeight))
	} else {
		sb.WriteString(fmt.Sprintf(` font-weight: %s;`, df.FontWeight))
	}
	if s.Decoration != "" {
		sb.WriteString(fmt.Sprintf(` text-decoration: %s;`, s.Decoration))
	} else {
		sb.WriteString(fmt.Sprintf(` text-decoration: %s;`, df.Decoration))
	}
	sb.WriteString(fmt.Sprintf(`">%s</span>`, elem))
	return Elem(sb.String())
}

type ImgElem struct {
	Alt    string
	Link   string
	Height int
	Width  int
}

func Img(src string, el *ImgElem) Elem {
	s := fmt.Sprintf(`<img src="%s"`, src)
	if el.Height != 0 {
		s += fmt.Sprintf(` height="%d"`, el.Height)
	}
	if el.Width != 0 {
		s += fmt.Sprintf(` width="%d"`, el.Width)
	}
	if el.Alt != "" {
		s += fmt.Sprintf(` alt="%s"`, el.Alt)
	}
	s += `style="display:block;margin-right:auto;margin-left:auto;width:100%; height:auto;"`
	// s += ` style="display:block;margin-right:auto;margin-left:auto;"`
	s += ">"
	elm := Elem(s)
	if el.Link != "" {
		elm = A(elm, el.Link)
	}
	return imgSection(elm)
}

func imgSection(elem Elem) Elem {
	s := fmt.Sprintf(`
	<table role="presentation" width="100%%" border="0" cellpadding="0" cellspacing="0" style="min-width:100%%">
	<tbody>
	   <tr style="white-space:nowrap;background-color:#ffffff">
		  <td align="center" style="background-color:#ffffff;padding-left:40px;padding-right:40px">
			 <div> 
				%s
			 </div>
		  </td>
	   </tr>
	</tbody>
 </table>
`,
		elem,
	)
	return Elem(s)
}

func mergeElems(elems ...Elem) Elem {
	var sb strings.Builder
	for _, e := range elems {
		sb.WriteString(" " + string(e))
	}
	return Elem(sb.String())
}
