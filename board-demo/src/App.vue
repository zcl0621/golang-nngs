<template>
  <div id="app">
    <div style="display: flex;justify-content: center;align-items: center;margin: 32px">
      <div style="display: inline-block; width: 300px;height: 600px">
        <div>
          <button @click="reload_board19">19路</button>
          <button @click="reload_board13">13路</button>
        </div>

        <div>
          <button @click="km7">KM 7.5</button>
          <button @click="km6">KM 6.5</button>
          <button @click="km0">KM 0</button>
        </div>

        <div>
          <textarea v-model="sgf" rows="10" cols="35" placeholder="输入sgf"></textarea>
          <button @click="load_sgf">确定加载</button>
        </div>
        <div>
          <button @click="close_ownership">关闭形势</button>
        </div>


        <p>棋盘尺寸: {{ size }}</p>
        <p>贴目: {{ km }}</p>
        <p>黑方：{{ b_score }}</p>
        <p>白方: {{ w_score }}</p>
        <p>scoreLead: {{ scoreLead }}</p>
        <p>scoreSelfplay: {{ scoreSelfplay }}</p>
        <p>scoreStdev: {{ scoreStdev }}</p>
        <p>颜色: {{ now_color === "B" ? '黑' : '白' }}</p>

      </div>

      <div
          style="
        display: inline-block;
        width: 600px;
        height: 600px;
        background: coral;
      "
          id="board"
          @click="click"
      ></div>

    </div>
    <div style="display: flex;justify-content: center;align-items: center;margin-top: 32px">
      <pre>{{ analysis_data }}</pre>
    </div>
  </div>
</template>

<script>
import "./wgo/wgo.min";
import "./wgo/wgo.player.min";
import axios from "axios";

export default {
  name: "App",
  components: {},
  data() {
    return {
      analysis_data: {},
      now_color: "B",
      b_score: 0,
      w_score: 0,
      scoreLead: 0,
      scoreSelfplay: 0,
      scoreStdev: 0,
      owners: [],
      player: "",
      sgf: "",
      size: 19,
      km: 7.5,
      lock: false,
      coordinates: {
        // draw on grid layer
        grid: {
          draw: function (args, board) {
            var ch, t, xright, xleft, ytop, ybottom;

            this.fillStyle = "rgba(0,0,0,1)";
            this.textBaseline = "middle";
            this.textAlign = "center";
            this.font = board.stoneRadius + "px " + (board.font || "");

            xright = board.getX(-0.75);
            xleft = board.getX(board.size - 0.25);
            ytop = board.getY(-0.75);
            ybottom = board.getY(board.size - 0.25);

            for (var i = 0; i < board.size; i++) {
              ch = i + "A".charCodeAt(0);
              if (ch >= "I".charCodeAt(0)) ch++;

              t = board.getY(i);
              this.fillText(board.size - i, xright, t);
              this.fillText(board.size - i, xleft, t);

              t = board.getX(i);
              this.fillText(String.fromCharCode(ch), t, ytop);
              this.fillText(String.fromCharCode(ch), t, ybottom);
            }

            this.fillStyle = "black";
          },
        },
      },
    };
  },
  methods: {
    to_obtain_coordinate: function (event) {
      // 捕捉落子的x y 坐标
      let board = this.player.board;
      let x, y;
      x = event.offsetX * board.pixelRatio;
      y = event.offsetY * board.pixelRatio;
      x -= board.left;
      x /= board.fieldWidth;
      x = Math.round(x);
      y -= board.top;
      y /= board.fieldHeight;
      y = Math.round(y);
      return {
        x: x >= board.size ? -1 : x,
        y: y >= board.size ? -1 : y,
      };
    },
    click(event) {
      this.player.board.removeObject(this.owners);
      let obj = this.to_obtain_coordinate(event);
      if (this.player.kifuReader.game.isValid(obj.x, obj.y) && !this.lock) {
        this.lock = true;
        this.player.kifuReader.node.appendChild(
            new WGo.KNode({
              move: {
                x: obj.x,
                y: obj.y,
                c: this.player.kifuReader.game.turn,
              },
            })
        );
        this.player.next(this.player.kifuReader.node.children.length - 1);
        this.sgf = this.player.kifuReader.kifu.toSgf().replace(/[\r\n]/g, "");
        this.getOwnerShip();
      }
    },
    make_player() {
      // 创建棋盘
      let elem = document.getElementById("board");
      this.player = new WGo.BasicPlayer(elem, {
        sgf: "(;SZ[19]KM[7.5])",
        layout: {
          left: "",
          bottom: "",
        },
        board: {
          section: {
            top: -0.5,
            right: -0.5,
            bottom: -0.5,
            left: -0.5
          },
        },
        enableWheel: false,
        markLastMove: false,
        displayVariations: false,
      });
      this.player.board.addCustomObject(this.coordinates);
    },
    load_sgf() {
      this.player.board.removeObject(this.owners);
      this.player.loadSgf(this.sgf);

      let p = WGo.clone(this.player.kifuReader.path);
      p.m += 1000;
      this.player.goTo(p);

      this.getOwnerShip();
    },
    getOwnerShip() {
      axios
          .post(
              "http://127.0.0.1:8080/api/katago-analysis/demo",
              {
                sgf: this.sgf,
              }
          )
          .then((res) => {
            this.analysis_data = res.data['analysis_data'];
            this.now_color = this.analysis_data['rootInfo']['currentPlayer'];
            this.scoreStdev = this.analysis_data['rootInfo']['scoreStdev'];
            this.scoreLead = this.analysis_data['rootInfo']['scoreLead'];
            this.scoreSelfplay = this.analysis_data['rootInfo']['scoreSelfplay'];
            let ownership = res.data['ownership']
            this.owners = [];
            this.b_score = 0;
            this.w_score = 0;
            let controversy_count = 0;
            for (let i = 0; i < ownership.length; i++) {
              let c;
              if (Math.abs(ownership[i].size) < 0.25) {
                if (ownership[i].c === 1) {
                  c = "pink";
                } else if (ownership[i].c === 2) {
                  c = "blue";
                } else {
                  c = "yellow";
                }
              } else {
                if (ownership[i].c === 1) {
                  c = "black";
                } else if (ownership[i].c === 2) {
                  c = "white";
                } else {
                  c = "yellow";
                }
              }
              let sq = {
                x: ownership[i].x,
                y: ownership[i].y,
                type: "SQ",
                c: c,
                size: Math.abs(ownership[i].size),
              };
              this.player.board.addObject(sq);
              this.owners.push(sq);
              let size = ownership[i].size;
              if (size >= 0.5) {
                this.b_score = this.b_score + 1;
              } else if (size < 0.5 && size >= 0.25) {
                this.b_score += 0.5;
              } else if (size <= -0.25 && size > -0.5) {
                this.w_score += 0.5;
              } else if (size <= -0.5) {
                this.w_score = this.w_score + 1;
              } else {
                controversy_count++;
              }
            }
            this.lock = false;
          })
          .catch((err) => {
            console.log(err);
            alert(err.toString())
            this.lock = false;
          });
    },
    close_ownership() {
      this.player.board.removeObject(this.owners);
    },
    reload_board19() {
      this.size = 19;
      this.reload();
    },
    reload_board13() {
      this.size = 13;
      this.reload();
    },
    km7() {
      this.km = 7.5;
      this.reload();
    },
    km6() {
      this.km = 6.5;
      this.reload();
    },
    km0() {
      this.km = 0;
      this.reload();
    },
    reload() {
      this.sgf = "";
      this.now_color = "B";
      this.b_score = 0;
      this.w_score = 0;
      this.scoreLead = 0;
      this.scoreSelfplay = 0;
      this.scoreStdev = 0;
      this.player.board.removeObject(this.owners);
      this.player.loadSgf(`(;SZ[${this.size}]KM[${this.km}])`);
    }
  },
  mounted() {
    this.make_player();
  },
};
</script>

<style lang="less">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
