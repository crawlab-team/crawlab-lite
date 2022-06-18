<script>
import * as Vue from 'vue'
export default function render(_props, _context) {
  const context = {
    ..._context,
    props: _props,
    data: _context.attr,
    children: _context.slots,
  }
  const { icon, title } = context.props
  const vnodes = []

  if (icon) {
    // vnodes.push(<svg-icon icon-class={icon}/>)
    const style = {
      'margin-right': '5px',
      'z-index': 999,
    }
    vnodes.push(<span class={icon} style={style} />)
  }

  if (title) {
    vnodes.push(
      <span class="title" slot="title">
        {title}
      </span>
    )
  }
  return vnodes
}
</script>
