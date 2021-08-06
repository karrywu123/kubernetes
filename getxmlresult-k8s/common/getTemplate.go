package common

func GetYmlTemplate(projectType string) (result string) {
	if projectType == "jar" || projectType == "target.zip" {
		result = `
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-service
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}-service
spec:
  type: NodePort
  selector:
    app: {{.ProjectName}}
  ports:
  {{ range .Ports }}
  - name: port-{{ . }}
    protocol: TCP
    port: {{ . }} #服务端口, 内部可访问
    targetPort: {{ . }} #Pod端口
  {{- end }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ProjectName}}
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
    spec:
      containers:
        - name: {{.ProjectName}}
          image: docker-harbor.clubs999.com/public/jboss/wildfly:10.0.0.Final
          command: ["/bin/sh"]
          args: ["-c", "/usr/lib/jvm/java/bin/java -Duser.dir=\"/data\" -Duser.timezone=GMT+8 -Xms1024m -Xmx10240m -jar /data/{{.Javaname}}"]
          securityContext:
            enabled: true
            privileged: true
            allowPrivilegeEscalation: true
            runAsUser: 0
            fsGroup: 0
          ports:
          {{ range .Ports }}
            - containerPort: {{ . }}
          {{- end }}
          volumeMounts:
            - mountPath: "/data"
              subPath: "{{.Namespace}}/{{.ProjectName}}" #文件路径,以挂载的pvc为基准
              name: {{.ProjectName}}-data
          env:
            - name: TZ
              value: "Asia/Shanghai"
            - name: LANG
              value: "en_US.UTF-8"
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "2048Mi"
              cpu: "2000m"
      volumes:
        - name: {{.ProjectName}}-data
          persistentVolumeClaim:
            claimName: gluster-claim-pub
		`
	}
	if projectType == "jboss" {
		result = `
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-service
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}-service
spec:
  type: NodePort
  selector:
    app: {{.ProjectName}}
  ports:
  {{ range .Ports }}
  - name: port-{{ . }}
    protocol: TCP
    port: {{ . }} #服务端口, 内部可访问
    targetPort: {{ . }} #Pod端口
  {{- end }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ProjectName}}
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
    spec:
      containers:
        - name: {{.ProjectName}}
          image: docker-harbor.clubs999.com/public/jboss/wildfly:10.0.0.Final
          command: ["/bin/sh"]
          args: ["-c", "/opt/jboss/wildfly/bin/standalone.sh -b 0.0.0.0 -bmanagement 0.0.0.0"]
          securityContext:
            enabled: true
            privileged: true
            allowPrivilegeEscalation: true
            runAsUser: 0
            fsGroup: 0
          ports:
          {{ range .Ports }}
          - containerPort: {{ . }}
          {{- end }}
          volumeMounts:
            - mountPath: "/opt/jboss/wildfly/standalone/deployments/"
              subPath: "{{.Namespace}}/{{.ProjectName}}" #文件路径,以挂载的pvc为基准
              name: {{.ProjectName}}-data
            - mountPath: "/opt/jboss/wildfly/bin/standalone.conf" #挂载jboss配置文件，可以自定义
              subPath: "standalone.conf"
              name: jboss-conf
          env:
            - name: TZ
              value: "Asia/Shanghai"
            - name: LANG
              value: "en_US.UTF-8"
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "2048Mi"
              cpu: "2000m"
      volumes:
        - name: {{.ProjectName}}-data
          persistentVolumeClaim:
            claimName: gluster-claim-pub
        - name: jboss-conf
          configMap:
            name: jboss-standalone-conf
		`
	}
	if projectType == "html" {
		result = `
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-service
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}-service
spec:
  type: NodePort
  selector:
    app: {{.ProjectName}}
  ports:
  {{ range .Ports }}
  - name: port-{{ . }}
    protocol: TCP
    port: {{ . }} #服务端口, 内部可访问
    targetPort: {{ . }} #Pod端口
  {{- end }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}
  namespace: {{.Namespace}}
  labels:
    app: {{.ProjectName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ProjectName}}
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
    spec:
      containers:
        - name: {{.ProjectName}}
          image: docker-harbor.clubs999.com/public/nginx:latest
          ports:
          {{ range .Ports }}
          - containerPort: {{ . }}
          {{- end }}
          volumeMounts:
            - mountPath: "/usr/share/nginx/html"
              subPath: "{{.Namespace}}/html/{{.ProjectName}}" #文件路径,以挂载的pvc为基准
              name: {{.ProjectName}}-data
          env:
            - name: TZ
              value: "Asia/Shanghai"
            - name: NGINX_PORT
              value: "80"
          resources:
            requests:
              memory: "256Mi"
              cpu: "50m"
            limits:
              memory: "2048Mi"
              cpu: "2000m"
      volumes:
        - name: {{.ProjectName}}-data
          persistentVolumeClaim:
            claimName: gluster-claim-pub
		`
	}

	return
}

func GetNginxTemplate(projectType string) (result string) {
	if projectType == "html" {
		result = `
upstream {{.ProjectName}} {
        {{ range .IpPorts }}
        server {{ . }} max_fails=2 fail_timeout=60s;
        {{- end }}
    }

server{
    listen 80;
    server_name {{ range .Domains }} {{ . }} {{- end }};
    server_tokens off;
    access_log  /var/log/nginx/{{.ProjectName}}.access.log main;
    add_header 'Access-Control-Allow-Origin' '*';
    add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
    add_header 'Access-Control-Allow-Credentials' 'true';
    location / {
        if ($http_user_agent ~* "WordPress") {
                return 502;
        }

        if ($http_user_agent ~* "spider") {
                return 502;
        }

        proxy_pass         http://{{.ProjectName}};
        proxy_set_header   Host             $host;
        proxy_set_header   X-Real-IP        $remote_addr;
        #proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-For  $http_x_forwarded_for;
        proxy_set_header   X-Forwarded-Scheme  "http";
        proxy_set_header   Remote-real-ip   $http_x_forwarded_for;
    }
        error_page  404 403        /40x.html;
        location = /40x.html {
                root   /usr/share/nginx/html;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
                root   /usr/share/nginx/html;
        }
}
`
	} else {
		result = `
upstream {{.ProjectName}} {
        {{ range .IpPorts }}
        server {{ . }} max_fails=2 fail_timeout=60s;
        {{- end }}
    }

server{
    listen 80;
    server_name {{ range .Domains }} {{ . }} {{- end }};
    server_tokens off;
    access_log  /var/log/nginx/{{.ProjectName}}.access.log main;
    location / {
        if ($http_user_agent ~* "WordPress") {
                return 502;
        }

        if ($http_user_agent ~* "spider") {
                return 502;
        }

        proxy_pass         http://{{.ProjectName}};
        proxy_set_header   Host             $host;
        proxy_set_header   X-Real-IP        $remote_addr;
        #proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-For  $http_x_forwarded_for;
        proxy_set_header   X-Forwarded-Scheme  "http";
        proxy_set_header   Remote-real-ip   $http_x_forwarded_for;
    }
        error_page  404 403        /40x.html;
        location = /40x.html {
                root   /usr/share/nginx/html;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
                root   /usr/share/nginx/html;
        }
}
		`
	}

	return
}
