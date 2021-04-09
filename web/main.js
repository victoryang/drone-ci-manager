const Home = { 
	template: '<div>home</div>',
	data: function(){
		return {
			projects: []
		}
	},
	created: function() {
		this.getProject()
	},
	watch: {
		"$route": 'getProject'
	},
	methods: {
		getProject() {
			this.$http.get("/projects/").then((response) => {
				console.log(response.data);
				this.projects = response.data;
			}, (response) => {
				alert("请求异常");
				this.$router.back();
			})
		},

		toCreateProject() {
			this.$prompt('Project name', '创建项目', {
				confirmButtonText: 'OK',
				cancelButtonText: 'Cancel',
			}).then(({ value }) => {
				this.$http.post("/projects/"+value).then((response) => {
					this.$message({
						type: 'success',
						message: value + ' 创建成功'
					});
					this.$router.push('/');
				}, (response) => {
					this.$message({
						type: 'info',
						message: '创建失败'
					});
				})
			}).catch(() => {
				this.$message({
					type: 'info',
					message: 'Canceled'
				});
			});
		},
		toHome() {
			this.$router.push('/');
		}
	}
}

const Project = {
	template: '<div>failed to load project page</div>',
	data() {
		return {
			buildInfo: {
				build_cmd: 'test',
				target: '',
				unzip_dir: '',
				lang: '',
				build_dependency: '',
				start_cmd: '',
				stop_cmd: '',
				pre_cmd: '',
				http_port: '',
				rpc_port: '',
			},
			ruleForm: {
				from_image: '',
				envs: [],
			},
			rules: {
				from_image: [
					{ required: true, message: 'Please select from image', trigger: 'change' }
				],
				envs: [
					{ type: 'array', required: true, message: 'Please select at least one env', trigger: 'change' }
				],
			}
		};
	},
	created() {
		this.fetchData();
	},
	watch: {
		'$route': 'fetchData'
	},
	methods: {
		fetchData () {
			this.$http.get("/projects/" + this.$route.params.project).then((response) => {
				console.log(response.data);
				this.buildInfo = response.data;
			}, (response) => {
				alert("请求异常");
				this.$router.back();
			})
		},
		createScrpit() {
			this.$http.post("/projects/" + this.$route.params.project + "/scripts", {
				param: {
					unzip_dir: this.buildInfo.unzip_dir,
					build_dependency: this.buildInfo.build_dependency,
					start_cmd: this.buildInfo.start_cmd,
					stop_cmd: this.buildInfo.stop_cmd,
					pre_cmd: this.buildInfo.pre_cmd,
					http_port: this.buildInfo.http_port,
					rpc_port: this.buildInfo.rpc_port,
					from_image: this.ruleForm.from_image,
				},
				envs: this.ruleForm.envs,
			}).then((response) => {
				this.$message("创建成功");
			}, (response) => {
				alert("创建失败: "+response.data.ErrMessage);
				this.$router.back();
			})
		},
	}
}

req = new XMLHttpRequest();
req.open('GET', 'home.html', false);
req.send(null);
Home.template = req.responseText;

req = new XMLHttpRequest();
req.open('GET', 'project.html', false);
req.send(null);
Project.template = req.responseText;

const router = new VueRouter({
	//mode: 'history',
	routes: [
		{ 
			path: '/', 
			component: Home,
			children: [{
				path: 'projects/:project',
				component: Project
			}]
		}
	]
})

new Vue({
	created() {
		this.$router.push('/')
	},
	router,
}).$mount('#app')