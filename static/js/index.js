  // 定义一个防抖函数
  function debounce(func, wait) {
    // 定义一个变量用来存储定时器的返回值
    let timeout;
    // 返回一个新的函数
    return function () {
      // 获取当前函数的执行上下文和参数
      let context = this;
      let args = arguments;
      // 如果已经存在定时器，就清除它
      if (timeout) {
        clearTimeout(timeout);
      }
      // 设置一个新的定时器，延迟执行原函数
      timeout = setTimeout(function () {
        func.apply(context, args);
      }, wait);
    };
  }

  const app = new Vue({
    el: '#app',
    data: {
      images: [],
    },
    computed: {
      imagesReverser() {
        return this.images.slice().reverse();
      }
    },
    methods: {
      debounceDraw: debounce(function (e) {
        console.log('draw');
        this.wf.draw();
      }, 300),
      imageLoad() {
        this.debounceDraw();
      },
      async fetchImages() {
        const res = await fetch(`/events?e=images&page=1`)
        if (res.status >= 300) {
          alert(await res.text());
          return
        }
        const data = await res.json()

        if (!data.success) {
          alert(data.message)
          return
        }

        this.images = data.data;
        Vue.nextTick(() => {
          this.debounceDraw();
        })
      },
      async del(item) {
        // 弹出一个确认框，显示"你确定要删除这个文件吗？"
        const result = confirm("你确定要删除这个图片吗？");
        if (!result) return;

        const res = await fetch(`/events?e=del&hash=${item.hash}`)
        if (res.status >= 300) {
          alert(await res.text());
          return
        }

        const data = await res.json()

        if (!data.success) {
          alert(data.message)
          return
        }

        const i = this.images.findIndex(e => e.hash == item.hash)
        this.images.splice(i, 1);
        alert('删除成功')

        Vue.nextTick(() => {
          this.debounceDraw();
        })
      },
      async uploadFile(e) {
        const file = e.target.files[0];

        // 创建一个 FormData 对象
        const formData = new FormData();

        // 添加文件数据
        formData.append("file", file);

        const res = await fetch("/events?e=upload", {
          method: "POST",
          body: formData,
        });

        if (res.status >= 300) {
          alert(await res.text());
          return
        }

        const data = await res.json()

        if (!data.success) {
          alert(data.message)
          return
        }

        this.images.push(data.data)
      },
      async copyUrl(item) {
        try {
          // 调用writeText方法，将文本写入剪贴板
          await navigator.clipboard.writeText(item.url);
          // 显示成功提示
          alert('内容已复制到剪贴板');
        } catch (err) {
          // 显示错误信息
          alert('复制失败: ' + err);
        }
      }
    },
    mounted() {
      this.wf = new window.waterfall.Waterfall({
        root: ".images",
        item: ".item",
        alignment: window.waterfall.WaterfallAlignment.center,
        reverse: true,
      });

      this.fetchImages();
    }
  })
