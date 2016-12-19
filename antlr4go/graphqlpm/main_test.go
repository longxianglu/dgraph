package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

var q1 = `
{
	al(xid: "alice") {
		status
		_xid_
		follows {
			status
			_xid_
			follows {
				status
				_xid_
				follows {
					_xid_
					status
				}
			}
		}
		status
		_xid_
	}
}
`

var q2 = `query queryName {
		me(uid : "0x0a") {
			friends {
				name
			}
			gender,age
			hometown
		}
	}
`

var q3 = `
{
  debug(xid: "m.0bxtg") {
    type.object.name.en
    film.actor.film {
      film.performance.film {
        film.film.directed_by {
          type.object.name.en
        }
      }
    }
  }
}
`

var q4 = `
{
  debug(_xid_: "m.06pj8") {
    type.object.name.en
    film.director.film {
      type.object.name.en
      film.film.initial_release_date
      film.film.country
      film.film.starring {
        film.performance.actor {
          type.object.name.en
        }
        film.performance.character {
          type.object.name.en
        }
      }
      film.film.genre {
        type.object.name.en
      }
    }
  }
}
`

func TestQueryParse(t *testing.T) {
	input := antlr.NewInputStream(q4)
	lexer := NewGraphQLPMLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewGraphQLPMParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Document()
	antlr.ParseTreeWalkerDefault.Walk(newMyListener(), tree)
}

func runParser(q string, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := antlr.NewInputStream(q)
		lexer := NewGraphQLPMLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := NewGraphQLPMParser(stream)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.BuildParseTrees = true
		// uptill here we have a cost of : 19000 for q1
		// next call makes it 100 times more costly to : 1800000
		_ = p.Document()
	}
}

func BenchmarkQuery(b *testing.B) {
	b.Run("q1", func(b *testing.B) { runParser(q1, b) })
	b.Run("q2", func(b *testing.B) { runParser(q2, b) })
	b.Run("q3", func(b *testing.B) { runParser(q3, b) })
	b.Run("q4", func(b *testing.B) { runParser(q4, b) })
}

// 1. Angelina jolie movies
// {
//   debug(_xid_: m.0f4vbz) {
//     type.object.name.en
//     film.actor.film {
//       film.performance.film {
//         type.object.name.en
//       }
//     }
//   }
// }
// first server run :
// "server_latency": {
//     "json": "528.087�s",
//     "parsing": "21.618164ms",
//     "processing": "26.746332ms",
//     "total": "48.898036ms"
//   }
// next run:
// "server_latency": {
//     "json": "474.959�s",
//     "parsing": "153.617�s",
//     "processing": "768.69�s",
//     "total": "1.401714ms"
//   }

// 2. movies directed by Steven Spielberg
// {
//   debug(_xid_: m.06pj8) {
//     type.object.name.en
//     film.director.film  {
//       film.film.genre {
//         type.object.name.en
//       }
//     }
//   }
// }
// first run :
// "server_latency": {
//     "json": "1.045405ms",
//     "parsing": "158.405�s",
//     "processing": "7.267112ms",
//     "total": "8.476755ms"
//   }
// next run:
//  "server_latency": {
//     "json": "1.050977ms",
//     "parsing": "216.808�s",
//     "processing": "914.389�s",
//     "total": "2.186685ms"
//   }
// 3. the movies acted by Brad Pitt
// {
//   debug(_xid_: m.0c6qh) {
//     type.object.name.en
//     film.actor.film {
//       film.performance.film {
//         type.object.name.en
//       }
//     }
//   }
// }
// first server run:
// "server_latency": {
//     "json": "874.075�s",
//     "parsing": "242.027�s",
//     "processing": "11.015496ms",
//     "total": "12.137608ms"
//   }
// next run:
// "server_latency": {
//     "json": "566.251�s",
//     "parsing": "159.576�s",
//     "processing": "988.121�s",
//     "total": "1.718198ms"
//   }
// 4. List of directors with whom Tom Hanks has worked
// {
//   debug(_xid_: m.0bxtg) {
//     type.object.name.en
//     film.actor.film {
//       film.performance.film {
//         film.film.directed_by {
//           type.object.name.en
//         }
//       }
//     }
//   }
// }
// first run:
// "server_latency": {
//     "json": "813.26�s",
//     "parsing": "205.994�s",
//     "processing": "17.100075ms",
//     "total": "18.122759ms"
//   }
// second run:
// "server_latency": {
//   "json": "879.414�s",
//   "parsing": "164.278�s",
//   "processing": "1.195895ms",
//   "total": "2.243513ms"
// }

// after running big query :
// dgraph was taking too much time;
// killed the server; restarted it..
// still taking too much time and memory.
// $ dgraph
// 2016/12/19 10:35:28 main.go:637: num_cpu: 4. prev_maxprocs: 4. Set max procs to num cpus
// Starting commit routine.
// 2016/12/19 10:35:29 main.go:685: grpc server started.
// 2016/12/19 10:35:29 main.go:686: http server started.
// 2016/12/19 10:35:29 main.go:687: Server listening on port 8080
// 2016/12/19 10:35:29 worker.go:80: Worker listening at address: [::]:12345
// NEW NODE GID, ID: [0, 1]
// NEW NODE GID, ID: [1, 1]
// Found hardstate: {Data:[] Metadata:{ConfState:{Nodes:[] XXX_unrecognized:[]} Index:0 Term:0 XXX_unrecognized:[]} XXX_unrecognized:[]}
// Found hardstate: {Data:[] Metadata:{ConfState:{Nodes:[] XXX_unrecognized:[]} Index:0 Term:0 XXX_unrecognized:[]} XXX_unrecognized:[]}
// Found 4 entries
// RESTARTING
// raft2016/12/19 10:35:30 INFO: 1 became follower at term 2
// raft2016/12/19 10:35:30 INFO: newRaft 1 [peers: [], term: 2, commit: 4, applied: 0, lastindex: 4, lastterm: 2]
// 2016/12/19 10:35:30 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// 2016/12/19 10:35:30 draft.go:310: group: 0 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:4}
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:3}
// ----------------------------
// ----------------------------
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:3}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:3}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:4}]
// raft2016/12/19 10:35:41 INFO: 1 is starting a new election at term 2
// raft2016/12/19 10:35:41 INFO: 1 became candidate at term 3
// raft2016/12/19 10:35:41 INFO: 1 received vote from 1 at term 3
// raft2016/12/19 10:35:41 INFO: 1 became leader at term 3
// raft2016/12/19 10:35:41 INFO: raft.node: 1 elected leader 1 at term 3
// 2016/12/19 10:35:41 draft.go:310: group: 1 Addr: "localhost:12345" leader: false dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:6}
// ----------------------------
// 2016/12/19 10:35:41 draft.go:310: group: 0 Addr: "localhost:12345" leader: false dead: false
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:3}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:6}]
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:7}
// ----------------------------
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:7}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:6}]
// Found 44288 entries
// RESTARTING
// raft2016/12/19 10:35:48 INFO: 1 became follower at term 2
// raft2016/12/19 10:35:48 INFO: newRaft 1 [peers: [], term: 2, commit: 44288, applied: 0, lastindex: 44288, lastterm: 2]
// raft2016/12/19 10:36:00 INFO: 1 is starting a new election at term 2
// raft2016/12/19 10:36:00 INFO: 1 became candidate at term 3
// raft2016/12/19 10:36:00 INFO: 1 received vote from 1 at term 3
// raft2016/12/19 10:36:00 INFO: 1 became leader at term 3
// raft2016/12/19 10:36:00 INFO: raft.node: 1 elected leader 1 at term 3
// 2016/12/19 10:36:25 draft.go:310: group: 0 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:8}
// ----------------------------
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:6}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:8}]

// 2016/12/19 10:38:04 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4207
// 2016/12/19 10:38:04 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:05 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:07 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3659
// 2016/12/19 10:38:14 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4302
// 2016/12/19 10:38:14 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:15 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:17 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3979
// 2016/12/19 10:38:19 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4157
// 2016/12/19 10:38:19 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:19 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:22 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4052
// 2016/12/19 10:38:24 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4237
// 2016/12/19 10:38:24 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:24 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:27 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4111
// 2016/12/19 10:38:28 draft.go:310: group: 0 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:13}
// ----------------------------
// 2016/12/19 10:38:28 draft.go:310: group: 0 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:9}
// ----------------------------
// 2016/12/19 10:38:28 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:10}
// ----------------------------
// 2016/12/19 10:38:28 draft.go:310: group: 0 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}
// ----------------------------
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:false RaftIdx:6}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:13}]
// 2016/12/19 10:38:28 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:12}
// ----------------------------
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:13}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:12}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:9}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:12}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:9}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:10}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:10}]
// 2016/12/19 10:38:29 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4283
// 2016/12/19 10:38:29 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:29 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:32 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4148
// 2016/12/19 10:38:34 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4320
// 2016/12/19 10:38:34 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:34 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:37 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4165
// 2016/12/19 10:38:39 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4310
// 2016/12/19 10:38:39 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:39 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:42 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4158
// 2016/12/19 10:38:44 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4315
// 2016/12/19 10:38:44 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:44 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:47 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4132
// 2016/12/19 10:38:49 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4299
// 2016/12/19 10:38:49 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:49 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:52 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4075
// 2016/12/19 10:38:54 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4250
// 2016/12/19 10:38:54 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:54 lists.go:94: Trying to free OS memory
// 2016/12/19 10:38:57 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4002
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:25}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:14}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:15}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:16}
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:17}
// ----------------------------
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:18}
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:25}]
// ----------------------------
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:19}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:20}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:21}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:22}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:23}
// ----------------------------
// 2016/12/19 10:38:58 draft.go:310: group: 1 Addr: "localhost:12345" leader: true dead: false
// ----------------------------
// ====== APPLYING MEMBERSHIP UPDATE: {NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:24}
// ----------------------------
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:18}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:15}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:16}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:17}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:19}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:20}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:21}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:22}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:23}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:24}]
// Group: 0. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:11}]
// Group: 1. List: [{NodeId:1 Addr:localhost:12345 Leader:true RaftIdx:14}]
// 2016/12/19 10:38:59 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4195
// 2016/12/19 10:38:59 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:38:59 lists.go:94: Trying to free OS memory
// 2016/12/19 10:39:02 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3903
// 2016/12/19 10:39:04 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4139
// 2016/12/19 10:39:04 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:39:04 lists.go:94: Trying to free OS memory
// 2016/12/19 10:39:07 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3770
// 2016/12/19 10:39:14 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4419
// 2016/12/19 10:39:14 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:39:15 lists.go:94: Trying to free OS memory
// 2016/12/19 10:39:16 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3422
// 2016/12/19 10:39:24 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4103
// 2016/12/19 10:39:24 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:39:25 lists.go:94: Trying to free OS memory
// 2016/12/19 10:39:26 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3203
// 2016/12/19 10:39:39 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4420
// 2016/12/19 10:39:39 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:39:40 lists.go:94: Trying to free OS memory
// 2016/12/19 10:39:41 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 2954
// 2016/12/19 10:42:29 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4143
// 2016/12/19 10:42:29 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:30 lists.go:94: Trying to free OS memory
// 2016/12/19 10:42:31 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3480
// 2016/12/19 10:42:39 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4198
// 2016/12/19 10:42:39 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:40 lists.go:94: Trying to free OS memory
// 2016/12/19 10:42:42 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 3873
// 2016/12/19 10:42:44 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4141
// 2016/12/19 10:42:44 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:44 lists.go:94: Trying to free OS memory
// 2016/12/19 10:42:47 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4009
// 2016/12/19 10:42:49 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4279
// 2016/12/19 10:42:49 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:50 lists.go:94: Trying to free OS memory
// 2016/12/19 10:42:52 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4127
// 2016/12/19 10:42:54 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4376
// 2016/12/19 10:42:54 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:54 lists.go:94: Trying to free OS memory
// 2016/12/19 10:42:57 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4226
// 2016/12/19 10:42:59 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4475
// 2016/12/19 10:42:59 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:42:59 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:02 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4322
// 2016/12/19 10:43:04 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4540
// 2016/12/19 10:43:04 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:04 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:07 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4426
// 2016/12/19 10:43:10 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4676
// 2016/12/19 10:43:10 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:10 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:12 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4533
// 2016/12/19 10:43:15 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4777
// 2016/12/19 10:43:15 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:15 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:18 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4607
// 2016/12/19 10:43:20 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4806
// 2016/12/19 10:43:20 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:20 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:23 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4668
// 2016/12/19 10:43:25 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4784
// 2016/12/19 10:43:25 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:25 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:27 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4709
// 2016/12/19 10:43:30 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4925
// 2016/12/19 10:43:30 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:30 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:33 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4779
// 2016/12/19 10:43:34 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4896
// 2016/12/19 10:43:34 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:35 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:38 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4822
// 2016/12/19 10:43:39 lists.go:89: Memory usage over threshold. STW. Allocated MB: 4989
// 2016/12/19 10:43:39 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:39 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:43 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4882
// 2016/12/19 10:43:44 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5039
// 2016/12/19 10:43:44 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:44 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:48 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4913
// 2016/12/19 10:43:49 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5053
// 2016/12/19 10:43:49 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:50 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:53 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4941
// 2016/12/19 10:43:54 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5068
// 2016/12/19 10:43:54 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:54 lists.go:94: Trying to free OS memory
// 2016/12/19 10:43:58 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4959
// 2016/12/19 10:43:59 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5087
// 2016/12/19 10:43:59 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:43:59 lists.go:94: Trying to free OS memory
// 2016/12/19 10:44:03 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4981
// 2016/12/19 10:44:05 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5146
// 2016/12/19 10:44:05 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:44:05 lists.go:94: Trying to free OS memory
// 2016/12/19 10:44:08 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 4996
// 2016/12/19 10:44:10 lists.go:89: Memory usage over threshold. STW. Allocated MB: 5143
// 2016/12/19 10:44:10 lists.go:91: Aggressive evict, committing to RocksDB
// 2016/12/19 10:44:10 lists.go:94: Trying to free OS memory
// 2016/12/19 10:44:13 lists.go:100: EVICT DONE! Memory usage after calling GC. Allocated MB: 5016
// ^C
// ashishnegi@ashish:~/work/golang/src/github.com/dgraph-io/dgraph/cmd/dgraphloader$
